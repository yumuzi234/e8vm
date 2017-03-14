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
	Open(p string) (*File, error)
}

// Output2 provides a simple output file system for storing build results.
type Output2 interface {
	// Create opens a file for writing.
	Create(p string) (io.WriteCloser, error)
}
