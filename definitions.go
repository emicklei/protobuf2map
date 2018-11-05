package protobuf2map

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	pp "github.com/emicklei/proto"
)

type Definitions struct {
	messages          map[string]*pp.Message
	enums             map[string]*pp.Enum
	filenamesRead     []string
	filenameToPackage map[string]string
}

func NewDefinitions() *Definitions {
	return &Definitions{
		messages:          map[string]*pp.Message{},
		enums:             map[string]*pp.Enum{},
		filenamesRead:     []string{},
		filenameToPackage: map[string]string{},
	}
}

// Read the proto definition from a filename.
func (d *Definitions) ReadFile(filename string) error {
	fileReader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	return d.ReadFrom(filename, fileReader)
}

// Read the proto definition from a Reader
// TODO: Recursively add all imports
func (d *Definitions) ReadFrom(filename string, reader io.Reader) error {
	for _, each := range d.filenamesRead {
		if each == filename {
			return nil
		}
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	d.filenamesRead = append(d.filenamesRead, filename)
	parser := pp.NewParser(bytes.NewReader(data))
	def, err := parser.Parse()
	if err != nil {
		return err
	}
	pkg := packageOf(def)
	d.filenameToPackage[filename] = pkg
	pp.Walk(def, pp.WithMessage(func(each *pp.Message) {
		d.AddMessage(pkg, each.Name, each)
	}))
	pp.Walk(def, pp.WithEnum(func(each *pp.Enum) {
		d.AddEnum(pkg, each.Name, each)
	}))
	return nil
}

// Returns the proto package name as declared in the proto filename.
func (d *Definitions) Package(filename string) (pkg string, ok bool) {
	pkg, ok = d.filenameToPackage[filename]
	return
}

func (d *Definitions) MessagesInPackage(pkg string) (list []*pp.Message) {
	for k, v := range d.messages {
		if strings.HasPrefix(k, pkg+".") {
			list = append(list, v)
		}
	}
	return
}

func (d *Definitions) Message(pkg string, name string) (m *pp.Message, ok bool) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	m, ok = d.messages[key]
	return
}

func (d *Definitions) Enum(pkg string, name string) (e *pp.Enum, ok bool) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	e, ok = d.enums[key]
	return
}

// AddEnum adds the Enum
func (d *Definitions) AddEnum(pkg string, name string, enu *pp.Enum) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.enums[key] = enu
}

func (d *Definitions) AddMessage(pkg string, name string, message *pp.Message) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.messages[key] = message
}

func packageOf(def *pp.Proto) string {
	for _, each := range def.Elements {
		if p, ok := each.(*pp.Package); ok {
			return p.Name
		}
	}
	return ""
}
