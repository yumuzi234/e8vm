package builds

import (
	"path/filepath"
)

// SingleFile creates a file map that can be used for single file compilation
func SingleFile(path string, o FileOpener) *FileSet {
	name := filepath.Base(path)
	file := &File{
		Name:   name,
		Path:   path,
		Opener: o,
	}
	ret := NewFileSet()
	ret.add(file)
	return ret
}

// SimplePkg creates a package that contains only one file
// and has no imports
func SimplePkg(p string, f string, o FileOpener) *PkgInfo {
	return &PkgInfo{
		Path:  p,
		Files: SingleFile(f, o),
	}
}
