package builds

import (
	"fmt"
)

type pkg struct {
	source  *source
	results *results
	path    string

	runTests bool

	srcMap  map[string]*File
	lang    Compiler
	files   []string
	imports map[string]*Import
	deps    []string

	pkg *Package
	err error
}

func newErrPkg(e error) *pkg { return &pkg{err: e} }

func newPkg(src *source, res *results, p string) *pkg {
	if !IsPkgPath(p) {
		return newErrPkg(fmt.Errorf("invalid path: %q", p))
	}

	lang := src.lang(p)
	if lang == nil {
		return newErrPkg(fmt.Errorf("invalid pacakge: %q", p))
	}

	srcMap, err := src.srcFileMap(p)
	if err != nil {
		return newErrPkg(err)
	}

	return &pkg{
		lang:    lang,
		source:  src,
		results: res,
		path:    p,
		srcMap:  srcMap,
		imports: make(map[string]*Import),
	}
}

func (p *pkg) fileSet() *FileSet {
	m := p.srcMap
	ret := NewFileSet()
	for _, f := range m {
		ret.add(f)
	}
	return ret
}

func (p *pkg) srcPackage() *SrcPackage {
	return &SrcPackage{
		Path:  p.path,
		Files: p.fileSet(),
	}
}
