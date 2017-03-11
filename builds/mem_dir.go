package builds

import (
	"sort"
)

type memDir struct {
	path  string
	files map[string]*memFile
}

func newMemDir(p string) *memDir {
	return &memDir{
		path:  p,
		files: make(map[string]*memFile),
	}
}

func (dir *memDir) list() []string {
	var ret []string
	for f := range dir.files {
		ret = append(ret, f)
	}
	sort.Strings(ret)
	return ret
}

func (dir *memDir) create(f string) *memFile {
	ret := newMemFile2()
	dir.files[f] = ret
	return ret
}

func (dir *memDir) open(f string) *memFile {
	return dir.files[f]
}

func (dir *memDir) file(f string) *memFile {
	ret, ok := dir.files[f]
	if !ok {
		return dir.create(f)
	}
	return ret
}
