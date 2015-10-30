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

// TODO(h8liu): this function is way too long; please refactor.
func (b *Builder) build(p string, forTest bool) (*pkg, []*lex8.Error) {
	ret := b.pkgs[p]
	if ret == nil {
		panic("build without prepare")
	}

	// already compiled
	if ret.compiled != nil {
		return ret, nil
	}

	if ret.buildStarted {
		e := fmt.Errorf("package %q circular depends itself", p)
		return ret, lex8.SingleErr(e)
	}

	ret.buildStarted = true
	lang := ret.lang

	for _, imp := range ret.imports {
		built, es := b.build(imp.Path, forTest)
		if es != nil {
			return nil, es
		}
		imp.Compiled = built.compiled
	}

	// ready to build this one
	if b.Verbose {
		fmt.Println(p)
	}

	// compile now
	pinfo := &PkgInfo{
		Path:   p,
		Src:    ret.srcMap(),
		Import: ret.imports,
		CreateLog: func(name string) io.WriteCloser {
			return b.home.CreateLog(p, name)
		},
	}
	compiled, es := lang.Compile(pinfo)
	if es != nil {
		return nil, es
	}
	ret.compiled = compiled

	lib := ret.compiled.Lib() // the linkable lib
	// a package with main entrance, build the bin

	main := ret.compiled.Main()
	if main != "" && lib.HasFunc(main) {
		log := lex8.NewErrorList()

		fout := b.home.CreateBin(p)
		lex8.LogError(log, b.link(lib, fout, main))
		lex8.LogError(log, fout.Close())

		if es := log.Errs(); es != nil {
			return nil, es
		}
	}

	if forTest {
		tests, testMain := ret.compiled.Tests()
		if testMain != "" && lib.HasFunc(testMain) {
			log := lex8.NewErrorList()
			if len(tests) > 0 {
				bs := new(bytes.Buffer)
				lex8.LogError(log, b.link(lib, bs, testMain))
				fout := b.home.CreateTestBin(p)

				img := bs.Bytes()
				_, err := fout.Write(img)
				lex8.LogError(log, err)
				lex8.LogError(log, fout.Close())
				if es := log.Errs(); es != nil {
					return nil, es
				}

				runTests(log, tests, img, b.Verbose)
				if es := log.Errs(); es != nil {
					return nil, es
				}
			}
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
