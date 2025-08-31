package gma

type Entry struct {
	name    string
	content []byte
	size    uint64
}

var HEADER = []byte("GMAD")

const VERSION = int8(3)

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Content() []byte {
	return e.content
}

func (e *Entry) Size() uint64 {
	return e.size
}
