package builds

import (
	"io"
)

// Input2 provides a simple input file system for building.
type Input2 interface {
	// HasDir checks if a directory exist.
	HasDir(p string) (bool, error)

	// ListDirs lists all the directories under a directory.
	ListDirs(p string) ([]string, error)

	// ListFiles lists all the files under a directory.
	ListFiles(p string) ([]string, error)

	// Open opens a file.
	Open(file string) *File
}

// Input provides input source files.
type Input interface {
	// HasPkg checks if a package exists.
	HasPkg(p string) bool

	// Pkgs lists all the packages with the particular prefix.
	Pkgs(prefix string) []string

	// Src lists the source files in a package.
	Src(path string) map[string]*File

	// Lang returns the language of a path.
	Lang(path string) *Lang
}

// Output provides writers for compiler outputs
type Output interface {
	// Output creates a compiler output, usually for debugging.
	Output(path, name string) io.WriteCloser

	// Bin creates the writer for generate the E8 binary image.
	Bin(path string) io.WriteCloser

	// TestBin creates the writer for generate the E8 binary image
	// for testing.
	TestBin(path string) io.WriteCloser
}

// Home is an interface that provides both input and output.
type Home interface {
	Input
	Output
}
