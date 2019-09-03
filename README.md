# protobuf2map

protobuf2map can decode a ProtocolBuffer marshaled message into a map[string]interface{} using its .proto definition.

This package was created to transform such messages directly to JSON without using the generated (Go) code to do the unmarshaling and marshaling. It can be used to inspect messages for quick view and debugging. Do not use this if performance is of importance.

	package main

	import (
			pd "github.com/emicklei/protobuf2map"
	)
	
	func main() {
		defs := pd.NewDefinitions()
		_ = defs.ReadFile("your.proto")
		protoBytes := []byte{} // your marshalled Protobuf message
		dec := pd.NewDecoder(defs, proto.NewBuffer(protoBytes))
		result, _ := dec.Decode("yours", "YourMessage")
		log.Printf("%#v",result)
	}
	

## how to compile the test

	protoc --go_out=. *.proto && mv test.pb.go test_pb_test.go