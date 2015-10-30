package build8

import (
	"io"

	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
	"e8vm.io/e8vm/sym8"
)

// File is a file in a package.
type File struct {
	Name string
	Path string
	io.ReadCloser
}

// Import is an import identity
type Import struct {
	Path     string
	Pos      *lex8.Pos
	Compiled Linkable
}

// Linkable is an interface for a linkable package
type Linkable interface {
	// Main is the entry point for linking. It will
	// be placed at the entry address of the image.
	Main() string

	// Tests are function symbols that should be preserved,
	// and sent into the image as an argument for running
	// testing.
	Tests() (tests map[string]uint32, main string)

	// Lib is the linkable object file.
	Lib() *link8.Pkg

	// Symbols returns all the top-level symbols in this package.
	// This is used for other packages to import and use this package.
	Symbols() (lang string, table *sym8.Table)
}

// Importer is an interface for importing required packages for compiling
type Importer interface {
	Import(name, path string, pos *lex8.Pos) // imports a package
}

// PkgInfo contains the information for compiling a package
type PkgInfo struct {
	Path   string
	Src    map[string]*File
	Import map[string]*Import

	CreateLog func(name string) io.WriteCloser
}

// Lang is a language compiler interface
type Lang interface {
	// IsSrc filters source file filenames
	IsSrc(filename string) bool

	// Prepare issues import requests
	Prepare(src map[string]*File, importer Importer) []*lex8.Error

	// Compile compiles a list of source files into a compiled linkable
	Compile(pinfo *PkgInfo) (Linkable, []*lex8.Error)
}
