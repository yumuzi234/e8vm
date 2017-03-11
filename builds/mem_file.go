package builds

import (
	"bytes"
)

type memFile struct {
	*bytes.Buffer
	path string
}

func newMemFile2() *memFile {
	return &memFile{
		Buffer: new(bytes.Buffer),
	}
}

func newMemFile(path string) *memFile {
	return &memFile{
		path:   path,
		Buffer: new(bytes.Buffer),
	}
}

func (f *memFile) Opener() FileOpener {
	bs := f.Buffer.Bytes()
	return BytesFile(bs)
}

func (f *memFile) Close() error { return nil }
