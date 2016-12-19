package builds

import (
	"io"
)

// SrcFile is a source file in a package
type SrcFile struct {
	Name string
	Path string // for printing compiler errors
	Open func() (io.ReadCloser, error)
}

// SrcPackage contains all the input of a package
type SrcPackage struct {
	Path  string     // package import path
	Lang  string     // language the package is written
	Hash  string     // a signature; empty for unknown
	Files []*SrcFile // list of source files
}

// SrcLoader loads packages.
type SrcLoader interface {
	LoadPkg(path string) *SrcPackage
}

// SrcLister lists packages of a given pattern.
type SrcLister interface {
	ListPkg(pattern string) []*SrcPackage
}
