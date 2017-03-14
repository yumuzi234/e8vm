package builds

import (
	"io"
)

// Input provides a simple input file system for building.
type Input interface {
	// HasDir checks if a directory exist.
	HasDir(p string) (bool, error)

	// ListDirs lists all the directories under a directory.
	ListDirs(p string) ([]string, error)

	// ListFiles lists all the files under a directory.
	ListFiles(p string) ([]string, error)

	// Open opens a file.
	Open(p string) (*File, error)
}

// Output provides a simple output file system for storing build results.
type Output interface {
	// Create opens a file for writing.
	Create(p string) (io.WriteCloser, error)
}
