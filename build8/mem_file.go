package build8

import (
	"bytes"
	"io"
	"io/ioutil"
)

type memFile struct {
	path string
	*bytes.Buffer
}

func newMemFile() *memFile {
	return &memFile{Buffer: new(bytes.Buffer)}
}

func (f *memFile) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(f.Buffer.Bytes()))
}

func (f *memFile) Close() error { return nil }
