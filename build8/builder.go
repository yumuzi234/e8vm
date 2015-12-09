package build8

import (
	"bytes"
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// Builder builds a bunch of packages.
type Builder struct {
	home Home
	pkgs map[string]*pkg

	Verbose bool
	InitPC  uint32
}

// NewBuilder creates a new builder with a particular home directory
func NewBuilder(home Home) *Builder {
	ret := new(Builder)
	ret.home = home
	ret.pkgs = make(map[string]*pkg)
	ret.InitPC = arch8.InitPC
	return ret
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

	return pkg, nil
}

func (b *Builder) link(p *link8.Pkg, out io.Writer, main string) error {
	job := link8.NewJob(p, main)
	job.InitPC = b.InitPC
	return job.Link(out)
}

func (b *Builder) buildImports(p *pkg, forTest bool) []*lex8.Error {
	for _, imp := range p.imports {
		built, es := b.build(imp.Path, forTest)
		if es != nil {
			return es
		}
		imp.Package = built.pkg
	}
	return nil
}

func (b *Builder) buildMain(p *pkg) []*lex8.Error {
	lib := p.pkg.Lib
	main := p.pkg.Main

	if main != "" && lib.HasFunc(main) {
		log := lex8.NewErrorList()

		fout := b.home.CreateBin(p.path)
		lex8.LogError(log, b.link(lib, fout, main))
		lex8.LogError(log, fout.Close())

		if es := log.Errs(); es != nil {
			return es
		}
	}

	return nil
}

func (b *Builder) runTests(p *pkg) []*lex8.Error {
	lib := p.pkg.Lib
	tests := p.pkg.Tests
	testMain := p.pkg.TestMain
	if testMain != "" && lib.HasFunc(testMain) {
		log := lex8.NewErrorList()
		if len(tests) > 0 {
			bs := new(bytes.Buffer)
			lex8.LogError(log, b.link(lib, bs, testMain))
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

func (b *Builder) build(p string, forTest bool) (*pkg, []*lex8.Error) {
	ret := b.pkgs[p]
	if ret == nil {
		panic("build without preparing")
	}

	// if already compiled, just return
	if ret.pkg != nil {
		return ret, nil
	}

	// circ dep check
	if ret.buildStarted {
		e := fmt.Errorf("package %q circular depends itself", p)
		return ret, lex8.SingleErr(e)
	}
	ret.buildStarted = true

	// build imports
	es := b.buildImports(ret, forTest)
	if es != nil {
		return nil, es
	}

	// report progress
	if b.Verbose {
		fmt.Println(p)
	}

	// compile
	compiled, es := ret.lang.Compile(b.makePkgInfo(ret))
	if es != nil {
		return nil, es
	}
	ret.pkg = compiled

	// build main
	es = b.buildMain(ret)
	if es != nil {
		return nil, es
	}

	// run tests
	if forTest {
		es := b.runTests(ret)
		if es != nil {
			return nil, es
		}
	}

	return ret, nil
}

// Build builds a package
func (b *Builder) Build(p string) []*lex8.Error {
	if _, es := b.prepare(p); es != nil {
		return es
	}

	_, es := b.build(p, false)
	return es
}

// BuildAll builds all packages, when andTest is also true,
// it will also test the packages.
func (b *Builder) BuildAll(andTest bool) []*lex8.Error {
	pkgs := b.home.Pkgs("")

	for _, p := range pkgs {
		if _, es := b.prepare(p); es != nil {
			return es
		}
	}

	for _, p := range pkgs {
		if _, es := b.build(p, andTest); es != nil {
			return es
		}
	}

	return nil
}
