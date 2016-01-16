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
	Path string
	Pos  *lex8.Pos
	*Package
}

// Package is an interface for a linkable package
type Package struct {
	// Lang is the language name this package used
	Lang string

	// Init is the init function of this package.
	// It is always a function that has no paramters.
	Init string

	// Main is the main entrance of this package, if any.
	Main string

	// TestMain is the main entrance for testing of this package, if any.
	TestMain string

	// Tests is the list of test cases, mapping from names to test ids.
	Tests map[string]uint32

	// Lib is the linkable library.
	Lib *link8.Pkg

	// Symbols stores all the symbols of this package.
	Symbols *sym8.Table
}

// Importer is an interface for importing required packages for compiling
type Importer interface {
	Import(name, path string, pos *lex8.Pos) // imports a package
}

// PkgSym is a pointer to a symbol in a package.
type PkgSym struct {
	Pkg, Sym string
}

// PkgInfo contains the information for compiling a package
type PkgInfo struct {
	Path   string
	Src    map[string]*File
	Import map[string]*Import

	Inits []*PkgSym

	// CreateLog creates the log file
	CreateLog func(name string) io.WriteCloser

	// AddFuncDebug adds debug information for a linking function.
	AddFuncDebug func(name string, pos *lex8.Pos, frameSize uint32)
}

// Lang is a language compiler interface
type Lang interface {
	// IsSrc filters source file filenames
	IsSrc(filename string) bool

	// Prepare issues import requests
	Prepare(src map[string]*File, importer Importer) []*lex8.Error

	// Compile compiles a list of source files into a compiled linkable
	Compile(pinfo *PkgInfo, opts *Options) (*Package, []*lex8.Error)
}
