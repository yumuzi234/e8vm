package builds

// MemPkg is a package in memory
type MemPkg struct {
	path  string
	files map[string]*memFile
	outs  map[string]*memFile
	bin   *memFile
	test  *memFile
	lib   *memFile
}

func newMemPkg(path string) *MemPkg {
	return &MemPkg{
		path:  path,
		outs:  make(map[string]*memFile),
		files: make(map[string]*memFile),
	}
}

// Path returns the path of the package
func (p *MemPkg) Path() string { return p.path }

// AddFile adds (or replaces) a source file in the package
func (p *MemPkg) AddFile(path, name, content string) {
	f := newMemFile()
	f.path = path
	f.WriteString(content)
	p.files[name] = f
}
