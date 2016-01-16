// Package asm8 prvoides an assembly language compiler for E8VM.
package asm8

import (
	"strings"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

type lang struct{}

func (lang) IsSrc(filename string) bool {
	return strings.HasSuffix(filename, ".s")
}

func (lang) Prepare(
	src map[string]*build8.File, imp build8.Importer,
) []*lex8.Error {
	if f := build8.OnlyFile(src); f != nil {
		return listImport(f.Path, f, imp)
	}

	f := src["import.s"]
	if f == nil {
		return nil
	}
	return listImport(f.Path, f, imp)
}

func buildSymTable(p *lib) *sym8.Table {
	t := sym8.NewTable()
	for _, sym := range p.symbols {
		if sym.Type == SymFunc || sym.Type == SymVar {
			t.Declare(sym)
		}
	}
	return t
}

func (lang) Compile(pinfo *build8.PkgInfo, opts *build8.Options) (
	*build8.Package, []*lex8.Error,
) {
	// resolve pass, will also parse the files
	pkg, es := resolvePkg(pinfo.Path, pinfo.Src)
	if es != nil {
		return nil, es
	}

	// import
	errs := lex8.NewErrorList()
	if pkg.imports != nil {
		for _, stmt := range pkg.imports.stmts {
			imp := pinfo.Import[stmt.as]
			if imp == nil || imp.Package == nil {
				errs.Errorf(stmt.Path.Pos, "import missing")
				continue
			}

			if imp.Lang != "asm8" {
				errs.Errorf(stmt.Path.Pos, "can only import asm8 package")
				continue
			}

			stmt.pkg = imp.Package
			if stmt.pkg == nil {
				panic("import missing")
			}
		}

		if es := errs.Errs(); es != nil {
			return nil, es
		}
	}

	// library building
	b := newBuilder(pinfo.Path)
	lib := buildLib(b, pkg)
	if es := b.Errs(); es != nil {
		return nil, es
	}

	ret := &build8.Package{
		Lang:    "asm8",
		Lib:     lib.Pkg,
		Main:    "main",
		Symbols: buildSymTable(lib),
	}
	return ret, nil
}

// Lang returns the assembly language builder for the building system
func Lang() build8.Lang { return lang{} }
