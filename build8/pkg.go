package build8

import (
	"fmt"

	"e8vm.io/e8vm/lexing"
)

type pkg struct {
	input  Input
	output Output
	path   string
	src    string

	runTests bool

	lang    Lang
	files   []string
	imports map[string]*Import
	deps    []string

	pkg *Package
	err error
}

func newErrPkg(e error) *pkg { return &pkg{err: e} }

func newPkg(in Input, out Output, p string) *pkg {
	if !isPkgPath(p) {
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

func (p *pkg) Import(name, path string, pos *lexing.Pos) {
	p.imports[name] = &Import{Path: path, Pos: pos}
}

var _ Importer = new(pkg)

/*
func (p *pkg) lastUpdate(suffix string) (*timeStamp, error) {
	ts := new(timeStamp)

	dirInfo, e := os.Stat(p.src)
	if e != nil {
		return nil, e
	}
	ts.update(dirInfo.ModTime())

	files, e := ioutil.ReadDir(p.src)
	if e != nil {
		return nil, e
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if p.lang.IsSrc(name) {
			ts.update(file.ModTime())
		}
	}

	return ts, nil
}

func (p *pkg) lastBuild() (*timeStamp, error) {
	ts := new(timeStamp)

	info, e := os.Stat(p.home.pkg(p.path))
	if !os.IsNotExist(e) {
		return nil, e
	}
	ts.update(info.ModTime())

	info, e = os.Stat(p.home.bin(p.path))
	if !os.IsNotExist(e) {
		return nil, e
	}
	ts.update(info.ModTime())

	return ts, nil
}
*/
