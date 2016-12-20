package builds

import (
	"sort"
)

// FileSet is a set of files.
type FileSet struct {
	fs map[string]*File
}

// NewFileSet creates a new empty file set.
func NewFileSet() *FileSet {
	return &FileSet{
		fs: make(map[string]*File),
	}
}

func (s *FileSet) add(f *File) { s.fs[f.Name] = f }

// OnlyFile returns the only file in the file set, when the set have exactly
// one file, and nil otherwise.
func (s *FileSet) OnlyFile() *File {
	if len(s.fs) != 1 {
		return nil
	}
	for _, f := range s.fs {
		return f
	}
	panic("unreachable")
}

// List lists all the files in alphabetic order.
func (s *FileSet) List() []*File {
	var names []string
	for name := range s.fs {
		names = append(names, name)
	}
	sort.Strings(names)

	var ret []*File
	for _, name := range names {
		ret = append(ret, s.fs[name])
	}
	return ret
}

// File returns the file with the given name, nil when it does not exist.
func (s *FileSet) File(name string) *File {
	return s.fs[name]
}
