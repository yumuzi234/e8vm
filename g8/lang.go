package g8

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type lang struct {
	golike bool
}

// Lang returns the G language builder for the building system
func Lang(golike bool) build8.Lang { return &lang{golike: golike} }

// LangGoLike returns the G lauguage that uses a subset of golang AST.
func LangGoLike() build8.Lang {
	return &lang{golike: true}
}

func (l *lang) IsSrc(filename string) bool {
	return strings.HasSuffix(filename, ".g")
}

func (l *lang) Prepare(
	src map[string]*build8.File, importer build8.Importer,
) []*lex8.Error {
	importer.Import("$", "asm/builtin", nil)
	if f := build8.OnlyFile(src); f != nil {
		return listImport(f.Path, f, importer, l.golike)
	}

	f := src["import.g"]
	if f == nil {
		return nil
	}
	return listImport(f.Path, f, importer, l.golike)
}

func makeBuilder(pinfo *build8.PkgInfo) *builder {
	b := newBuilder(pinfo.Path)
	initBuilder(b, pinfo.Import)
	return b
}

func initBuilder(b *builder, imp map[string]*build8.Import) {
	b.exprFunc = buildExpr2
	b.stmtFunc = buildStmt

	builtin, ok := imp["$"]
	if !ok {
		b.Errorf(nil, "builtin import missing for %q", b.path)
		return
	}

	declareBuiltin(b, builtin.Lib)
}

// parse all files
func (l *lang) parsePkg(pinfo *build8.PkgInfo) (
	map[string]*ast.File, []*lex8.Error,
) {
	var parseErrs []*lex8.Error
	asts := make(map[string]*ast.File)
	for name, src := range pinfo.Src {
		if filepath.Base(src.Path) != name {
			panic("basename in path is different from the file name")
		}

		f, rec, es := parse.File(src.Path, src, l.golike)
		if es != nil {
			parseErrs = append(parseErrs, es...)
		}
		if err := src.Close(); err != nil {
			parseErrs = append(parseErrs, &lex8.Error{Err: err})
		}

		if pinfo.ParseOutput != nil {
			pinfo.ParseOutput(name, rec.Tokens())
		}
		asts[name] = f
	}
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}

	return asts, nil
}

func output(w io.WriteCloser, f func(w io.Writer) error) error {
	err := f(w)
	if err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func outputIr(pinfo *build8.PkgInfo, b *builder) error {
	return output(pinfo.Output("ir"), func(w io.Writer) error {
		return ir.PrintPkg(w, b.p)
	})
}

func outputDeps(pinfo *build8.PkgInfo, g *dagvis.Graph) error {

	bs, err := json.MarshalIndent(g.Nodes, "", "    ")
	if err != nil {
		panic(err)
	}

	return output(pinfo.Output("deps"), func(w io.Writer) error {
		_, err := w.Write(bs)
		return err
	})
}

func outputDepMap(pinfo *build8.PkgInfo, deps []byte) error {
	return output(pinfo.Output("depmap"), func(w io.Writer) error {
		_, err := w.Write(deps)
		return err
	})
}

func (l *lang) outputDeps(pinfo *build8.PkgInfo, p *pkg) error {
	g := p.deps
	g, err := g.Rename(func(name string) (string, error) {
		if strings.HasSuffix(name, ".g") {
			return strings.TrimSuffix(name, ".g"), nil
		}
		return name, fmt.Errorf("filename suffix missing: %q", name)
	})

	if err != nil {
		return err
	}

	if err := outputDeps(pinfo, g); err != nil {
		return err
	}

	bs, err := dagvis.LayoutJSON(g.Reverse())
	if err != nil {
		return err
	}
	if err := outputDepMap(pinfo, bs); err != nil {
		return err
	}

	return err
}

func (l *lang) Compile(pinfo *build8.PkgInfo) (
	*build8.Package, []*lex8.Error,
) {
	ret := &build8.Package{
		Lang:     "g8",
		Init:     initName,
		Main:     startName,
		TestMain: testStartName,
	}

	// parsing
	asts, es := l.parsePkg(pinfo)
	if es != nil {
		return nil, es
	}

	// building
	b := makeBuilder(pinfo)
	if es = b.Errs(); es != nil {
		return nil, es
	}

	p := newPkg(asts)
	if es := p.build(b, pinfo); es != nil {
		return nil, es
	}

	// test mapping
	tests := make(map[string]uint32)
	for i, name := range p.testNames {
		tests[name] = uint32(i)
	}
	ret.Tests = tests

	// check deps
	if err := l.outputDeps(pinfo, p); err != nil {
		return nil, lex8.SingleErr(err)
	}

	ret.Symbols = p.tops
	if pinfo.Flags.StaticOnly { // static analysis stops here
		return ret, nil
	}

	var errs []*lex8.Error
	if ret.Lib, errs = ir.BuildPkg(b.p); errs != nil {
		return nil, errs
	}

	// add debug symbols
	//
	// Functions positionings only available after building.
	ir.AddDebug(b.p, pinfo.AddFuncDebug)

	// IR logging
	if err := outputIr(pinfo, b); err != nil {
		return nil, lex8.SingleErr(err)
	}

	return ret, nil
}
