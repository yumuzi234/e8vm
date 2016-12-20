package builds

import (
	"fmt"
)

type pkg struct {
	input  Input
	output Output
	path   string
	src    string

	runTests bool

	lang    Compiler
	files   []string
	imports map[string]*Import
	deps    []string

	pkg *Package
	err error
}

func newErrPkg(e error) *pkg { return &pkg{err: e} }

func newPkg(in Input, out Output, p string) *pkg {
	if !IsPkgPath(p) {
		return newErrPkg(fmt.Errorf("invalid path: %q", p))
	}

	lang := in.Lang(p)
	if lang == nil {
		return newErrPkg(fmt.Errorf("invalid pacakge: %q", p))
	} else if !in.HasPkg(p) {
		return newErrPkg(fmt.Errorf("package not found: %q", p))
	}

	return &pkg{
		lang:    lang,
		input:   in,
		output:  out,
		path:    p,
		imports: make(map[string]*Import),
	}
}

func (p *pkg) srcMap() map[string]*File { return p.input.Src(p.path) }

func (p *pkg) srcPackage() *SrcPackage {
	return &SrcPackage{
		Path:  p.path,
		Lang:  "", // TODO: specify language.
		Hash:  "",
		Files: p.input.Src(p.path),
	}
}
