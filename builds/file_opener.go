package builds

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// FileOpener opens a particular file for reading.
type FileOpener interface {
	Open() (io.ReadCloser, error)
}

type bytesFile struct{ bs []byte }

func (f *bytesFile) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(f.bs)), nil
}

// BytesFile creates a new opener that can open a set of bytes.
func BytesFile(bs []byte) FileOpener { return &bytesFile{bs} }

type strFile struct{ s string }

func (f *strFile) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(f.s)), nil
}

// StrFile creates a new opener that can open a string
func StrFile(s string) FileOpener { return &strFile{s} }

type pathFile struct{ p string }

func (f *pathFile) Open() (io.ReadCloser, error) {
	r, err := os.Open(f.p)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PathFile creates a new opener that opens a file on the file system.
func PathFile(p string) FileOpener { return &pathFile{p} }
