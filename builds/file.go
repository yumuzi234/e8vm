package builds

import (
	"io"
)

// File is a source file in a package
type File struct {
	Name   string // for sending into the complier
	Path   string // for printing compiler errors
	Opener FileOpener
}

// Open opens the file.
func (f *File) Open() (io.ReadCloser, error) {
	return f.Opener.Open()
}
