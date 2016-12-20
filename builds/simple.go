package builds

import (
	"path/filepath"
)

// SingleFile creates a file map that can be used for single file compilation
func SingleFile(path string, o FileOpener) map[string]*File {
	name := filepath.Base(path)
	file := &File{
		Name:   name,
		Path:   path,
		Opener: o,
	}
	return map[string]*File{name: file}
}

// SimplePkg creates a package that contains only one file
// and has no imports
func SimplePkg(p string, f string, o FileOpener) *PkgInfo {
	single := SingleFile(f, o)
	return &PkgInfo{
		Path: p,
		Src:  single,
	}
}
