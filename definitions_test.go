package protobuf2map

import "testing"

func TestAddFromFile(t *testing.T) {
	d := NewDefinitions()
	d.ReadFile("test.proto")
	if got, want := len(d.filenamesRead), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	m, ok := d.Message("protobuf2map", "Test")
	if !ok {
		t.Fail()
	}
	if got, want := m.Name, "Test"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
