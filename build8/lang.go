package build8

import (
	"io"

	"e8vm.io/e8vm/lexing"
	link8 "e8vm.io/e8vm/link"
	"e8vm.io/e8vm/syms"
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
	Pos  *lexing.Pos
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

	// Symbols stores all the symbols of this package.
	Symbols *syms.Table

	// Lib is the linkable library.
	Lib *link8.Pkg
}

// Importer is an interface for importing required packages for compiling
type Importer interface {
	Import(name, path string, pos *lexing.Pos) // imports a package
}

// Flags contains the flags for compiling a package
type Flags struct {
	StaticOnly bool // only perform static analysis
}

// PkgInfo contains the information for compiling a package
type PkgInfo struct {
	Path   string
	Src    map[string]*File
	Import map[string]*Import
	Flags  *Flags

	// Output creates an output file for the package.
	Output func(name string) io.WriteCloser

	// ParseOutput saves all the tokens of a file.
	ParseOutput func(file string, tokens []*lexing.Token)

	// AddFuncDebug adds debug information for a linking function.
	AddFuncDebug func(name string, pos *lexing.Pos, frameSize uint32)
}

// Lang is a language compiler interface
type Lang interface {
	// IsSrc filters source file filenames
	IsSrc(filename string) bool

	// Prepare issues import requests
	Prepare(src map[string]*File, importer Importer) []*lexing.Error

	// Compile compiles a list of source files into a compiled linkable
	Compile(pinfo *PkgInfo) (*Package, []*lexing.Error)
}
