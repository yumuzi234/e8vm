package build8

import (
	"bytes"
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// Builder builds a bunch of packages.
type Builder struct {
	home Home
	pkgs map[string]*pkg
	deps map[string][]string

	linkPkgs map[string]*link8.Pkg

	Verbose  bool
	InitPC   uint32
	RunTests bool
}

// NewBuilder creates a new builder with a particular home directory
func NewBuilder(home Home) *Builder {
	return &Builder{
		home:     home,
		pkgs:     make(map[string]*pkg),
		deps:     make(map[string][]string),
		linkPkgs: make(map[string]*link8.Pkg),
		InitPC:   arch8.InitPC,
	}
}

func (b *Builder) prepare(p string) (*pkg, []*lex8.Error) {
	saved := b.pkgs[p]
	if saved != nil {
		return saved, nil // already prepared
	}

	pkg := newPkg(b.home, p)
	b.pkgs[p] = pkg
	if pkg.err != nil {
		return pkg, nil
	}

	es := pkg.lang.Prepare(pkg.srcMap(), pkg)
	if es != nil {
		return pkg, es
	}

	for _, imp := range pkg.imports {
		impPkg, es := b.prepare(imp.Path)
		if es != nil {
			return pkg, es
		}

		if impPkg.err != nil {
			return pkg, []*lex8.Error{{
				Pos: imp.Pos,
				Err: impPkg.err,
			}}
		}
	}

	// mount the deps
	deps := make([]string, 0, len(pkg.imports))
	for _, imp := range pkg.imports {
		deps = append(deps, imp.Path)
	}
	b.deps[p] = deps

	return pkg, nil
}

func (b *Builder) link(out io.Writer, p *pkg, main string) error {
	var funcs []*link8.PkgSym

	addInit := func(p *pkg) {
		name := p.pkg.Init
		if name != "" && p.pkg.Lib.HasFunc(name) {
			funcs = append(funcs, &link8.PkgSym{p.path, name})
		}
	}

	for _, dep := range p.deps {
		addInit(b.pkgs[dep])
	}
	addInit(p)
	funcs = append(funcs, &link8.PkgSym{p.path, main})

	job := link8.NewJob(b.linkPkgs, funcs)
	job.InitPC = b.InitPC
	return job.Link(out)
}

func (b *Builder) fillImports(p *pkg) {
	for _, imp := range p.imports {
		imp.Package = b.pkgs[imp.Path].pkg
		if imp.Package == nil {
			panic("bug")
		}
	}
}

func (b *Builder) buildMain(p *pkg) []*lex8.Error {
	lib := p.pkg.Lib
	main := p.pkg.Main

	if main == "" || !lib.HasFunc(main) {
		return nil
	}

	log := lex8.NewErrorList()

	fout := b.home.CreateBin(p.path)
	lex8.LogError(log, b.link(fout, p, main))
	lex8.LogError(log, fout.Close())

	return log.Errs()
}

func (b *Builder) runTests(p *pkg) []*lex8.Error {
	lib := p.pkg.Lib
	tests := p.pkg.Tests
	testMain := p.pkg.TestMain
	if testMain != "" && lib.HasFunc(testMain) {
		log := lex8.NewErrorList()
		if len(tests) > 0 {
			bs := new(bytes.Buffer)
			lex8.LogError(log, b.link(bs, p, testMain))
			fout := b.home.CreateTestBin(p.path)

			img := bs.Bytes()
			_, err := fout.Write(img)
			lex8.LogError(log, err)
			lex8.LogError(log, fout.Close())
			if es := log.Errs(); es != nil {
				return es
			}

			runTests(log, tests, img, b.Verbose)
			if es := log.Errs(); es != nil {
				return es
			}
		}
	}

	return nil
}

func (b *Builder) makePkgInfo(p *pkg) *PkgInfo {
	return &PkgInfo{
		Path:   p.path,
		Src:    p.srcMap(),
		Import: p.imports,
		CreateLog: func(name string) io.WriteCloser {
			return b.home.CreateLog(p.path, name)
		},
	}
}

func (b *Builder) build(pkg *pkg) []*lex8.Error {
	b.fillImports(pkg)

	compiled, es := pkg.lang.Compile(b.makePkgInfo(pkg))
	if es != nil {
		return es
	}
	pkg.pkg = compiled
	b.linkPkgs[pkg.path] = pkg.pkg.Lib // add for linking

	if es := b.buildMain(pkg); es != nil {
		return es
	}

	if b.RunTests {
		if es := b.runTests(pkg); es != nil {
			return es
		}
	}

	return nil
}

func deps(node *dagvis.MapNode) []string {
	depNodes := dagvis.AllInsSorted(node)
	ret := make([]string, 0, len(depNodes))
	for _, dep := range depNodes {
		ret = append(ret, dep.Name)
	}
	return ret
}

// BuildPkgs builds a list of packages
func (b *Builder) BuildPkgs(pkgs []string) []*lex8.Error {
	for _, p := range pkgs {
		if _, es := b.prepare(p); es != nil {
			return es
		}
	}

	g := &dagvis.Graph{b.deps}
	g = g.Reverse()

	m, err := dagvis.Layout(g)
	if err != nil {
		return lex8.SingleErr(err)
	}
	// TODO: save the package dep map

	nodes := m.SortedNodes()
	for _, node := range nodes {
		if b.Verbose { // report progress
			fmt.Println(node.Name)
		}

		pkg := b.pkgs[node.Name]
		if pkg == nil {
			panic("package not prepared")
		}

		pkg.deps = deps(node)
		if es := b.build(pkg); es != nil {
			return es
		}
	}

	return nil
}

// Build builds a package.
func (b *Builder) Build(p string) []*lex8.Error {
	return b.BuildPkgs([]string{p})
}

// BuildAll builds all packages, when andTest is also true,
// it will also test the packages.
func (b *Builder) BuildAll() []*lex8.Error {
	return b.BuildPkgs(b.home.Pkgs(""))
}
