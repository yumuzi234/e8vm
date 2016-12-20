package builds

import (
	"bytes"
)

type memFile struct {
	path string
	*bytes.Buffer
}

func newMemFile() *memFile {
	return &memFile{Buffer: new(bytes.Buffer)}
}

func (f *memFile) Opener() FileOpener {
	bs := f.Buffer.Bytes()
	return BytesFile(bs)
}

func (f *memFile) Close() error { return nil }
