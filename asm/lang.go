// Package asm prvoides an assembly language compiler for the virtual
// instruction set.
package asm

import (
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/syms"
)

type lang struct{}

func (lang) Prepare(src *builds.FileSet) (
	*builds.ImportList, []*lexing.Error,
) {
	ret := builds.NewImportList()
	if f := src.OnlyFile(); f != nil {
		rc, err := f.Open()
		if err != nil {
			return nil, lexing.SingleErr(err)
		}
		if errs := listImport(f.Path, rc, ret); errs != nil {
			return nil, errs
		}
		return ret, nil
	}

	f := src.File("import.s")
	if f == nil {
		return nil, nil
	}
	rc, err := f.Open()
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
	if errs := listImport(f.Path, rc, ret); errs != nil {
		return nil, errs
	}
	return ret, nil
}

func buildSymTable(p *lib) *syms.Table {
	t := syms.NewTable()
	for _, sym := range p.symbols {
		if sym.Type == SymFunc || sym.Type == SymVar {
			t.Declare(sym)
		}
	}
	return t
}

func (lang) Compile(pinfo *builds.PkgInfo) (
	*builds.Package, []*lexing.Error,
) {
	// resolve pass, will also parse the files
	pkg, es := resolvePkg(pinfo.Path, pinfo.Files)
	if es != nil {
		return nil, es
	}

	// import
	errs := lexing.NewErrorList()
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

	ret := &builds.Package{
		Lang:    "asm8",
		Lib:     lib.Pkg,
		Main:    "main",
		Symbols: buildSymTable(lib),
	}
	return ret, nil
}

// Lang returns the assembly language builder for the building system
func Lang() *builds.Lang {
	return &builds.Lang{
		Ext:      "asm",
		Compiler: lang{},
	}
}
