package builds

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

// Add adds a file into the fileset.
func (s *FileSet) Add(f *File) { s.fs[f.Name] = f }

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

// File returns the file with the given name, nil when it does not exist.
func (s *FileSet) File(name string) *File {
	return s.fs[name]
}
