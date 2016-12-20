package builds

import (
	"io"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/link"
	"shanhu.io/smlvm/syms"
)

// Import is an import identity
type Import struct {
	Path string
	Pos  *lexing.Pos
	*Package
}

// Package is an interface for a linkable package
type Package struct {
	Lang    string      // the language name this package used
	Lib     *link.Pkg   // linkable object.
	Symbols *syms.Table // all the symbols

	Init     string            // the init function; always has no parameters.
	Main     string            // main entrance. optional.
	TestMain string            // test main entrance. optional.
	Tests    map[string]uint32 // list of tests. map from names to ids.
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

// Compiler is a language compiler that compiles a source package.
type Compiler interface {
	// Prepare issues import requests
	Prepare(src *SrcPackage) (*ImportList, []*lexing.Error)

	// Compile compiles a list of source files into a compiled linkable
	Compile(pinfo *PkgInfo) (*Package, []*lexing.Error)
}

// Lang contains the language info.
type Lang struct {
	Dir string
	Ext string
	Compiler
}
