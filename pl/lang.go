package pl

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/dagvis"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/parse"
)

type lang struct {
	golike bool
}

// Lang returns the G language builder for the building system
func Lang(golike bool) builds.Lang { return &lang{golike: golike} }

// LangGoLike returns the G lauguage that uses a subset of golang AST.
func LangGoLike() builds.Lang {
	return &lang{golike: true}
}

func (l *lang) IsSrc(filename string) bool {
	return strings.HasSuffix(filename, ".g")
}

func (l *lang) Prepare(src *builds.SrcPackage) (
	*builds.ImportList, []*lexing.Error,
) {
	ret := builds.NewImportList()
	ret.Add("$", "asm/builtin", nil)

	if f := builds.OnlyFile(src.Files); f != nil {
		if errs := listImport(f.Path, f, l.golike, ret); errs != nil {
			return nil, errs
		}
	}

	if src.Files != nil {
		f := src.Files["import.g"]
		if f == nil {
			return ret, nil
		}
		if errs := listImport(f.Path, f, l.golike, ret); errs != nil {
			return nil, errs
		}
	}
	return ret, nil
}

func makeBuilder(pinfo *builds.PkgInfo) *builder {
	b := newBuilder(pinfo.Path)
	initBuilder(b, pinfo.Import)
	return b
}

func initBuilder(b *builder, imp map[string]*builds.Import) {
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
func (l *lang) parsePkg(pinfo *builds.PkgInfo) (
	map[string]*ast.File, []*lexing.Error,
) {
	var parseErrs []*lexing.Error
	asts := make(map[string]*ast.File)
	for name, src := range pinfo.Src {
		if filepath.Base(src.Path) != name {
			panic("basename in path is different from the file name")
		}

		rc, err := src.Open()
		if err != nil {
			parseErrs = append(parseErrs, &lexing.Error{Err: err})
			continue
		}

		f, rec, es := parse.File(src.Path, rc, l.golike)
		if es != nil {
			parseErrs = append(parseErrs, es...)
		}
		if err := rc.Close(); err != nil {
			parseErrs = append(parseErrs, &lexing.Error{Err: err})
		}

		if pinfo.ParseOutput != nil && rec != nil {
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

func outputIr(pinfo *builds.PkgInfo, b *builder) error {
	return output(pinfo.Output("ir"), func(w io.Writer) error {
		return codegen.PrintPkg(w, b.p)
	})
}

func outputDeps(pinfo *builds.PkgInfo, g *dagvis.Graph) error {
	bs, err := json.MarshalIndent(g.Nodes, "", "    ")
	if err != nil {
		panic(err)
	}

	return output(pinfo.Output("deps"), func(w io.Writer) error {
		_, err := w.Write(bs)
		return err
	})
}

func outputDepMap(pinfo *builds.PkgInfo, deps []byte) error {
	return output(pinfo.Output("depmap"), func(w io.Writer) error {
		_, err := w.Write(deps)
		return err
	})
}

func (l *lang) outputDeps(pinfo *builds.PkgInfo, p *pkg) error {
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

func (l *lang) Compile(pinfo *builds.PkgInfo) (
	*builds.Package, []*lexing.Error,
) {
	ret := &builds.Package{
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
		return nil, lexing.SingleErr(err)
	}

	ret.Symbols = p.tops
	if pinfo.Flags.StaticOnly { // static analysis stops here
		return ret, nil
	}

	var errs []*lexing.Error
	if ret.Lib, errs = codegen.BuildPkg(b.p); errs != nil {
		return nil, errs
	}

	// add debug symbols
	// Functions positionings only available after building.
	codegen.AddDebug(b.p, pinfo.AddFuncDebug)

	// IR logging
	if err := outputIr(pinfo, b); err != nil {
		return nil, lexing.SingleErr(err)
	}

	return ret, nil
}
