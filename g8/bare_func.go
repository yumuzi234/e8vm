package g8

import (
	"fmt"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/g8/sempass"
	"e8vm.io/e8vm/lex8"
)

// because bare function also uses builtin functions that comes from the
// building system, we also need to make it a simple language: a
// language with only one (implicit) main function
// In fact, we can simple "inherit" the basic lang
type bareFunc struct{ *lang }

// BareFunc is a language where it only contains an implicit main function.
func BareFunc() build8.Lang { return bareFunc{new(lang)} }

func (bareFunc) Prepare(
	src map[string]*build8.File, importer build8.Importer,
) []*lex8.Error {
	importer.Import("$", "asm/builtin", nil)
	return nil
}

func buildBareFunc(b *builder, stmts []ast.Stmt) []*lex8.Error {
	tstmts, errs := sempass.BuildBareFunc(b.scope, stmts)
	if errs != nil {
		return errs
	}

	b.f = b.p.NewFunc(":start", nil, ir.VoidFuncSig)
	b.fretRef = nil
	b.b = b.f.NewBlock(nil)

	for _, stmt := range tstmts {
		buildStmt(b, stmt)
	}
	return nil
}

func findTheFile(pinfo *build8.PkgInfo) (*build8.File, error) {
	if len(pinfo.Src) == 0 {
		panic("no source file")
	} else if len(pinfo.Src) > 1 {
		return nil, fmt.Errorf("bare func %q has too many files", pinfo.Path)
	}

	for _, r := range pinfo.Src {
		return r, nil
	}
	panic("unreachable")
}

func (bare bareFunc) Compile(pinfo *build8.PkgInfo) (
	pkg *build8.Package, es []*lex8.Error,
) {
	// parsing
	theFile, e := findTheFile(pinfo)
	if e != nil {
		return nil, lex8.SingleErr(e)
	}
	stmts, es := parse.Stmts(theFile.Path, theFile)
	if es != nil {
		return nil, es
	}

	// building
	b := makeBuilder(pinfo, false)
	if es = b.Errs(); es != nil {
		return nil, es
	}
	if es := buildBareFunc(b, stmts); es != nil {
		return nil, es
	}
	if es := b.Errs(); es != nil {
		return nil, es
	}

	// logging
	if e := logIr(pinfo, b); e != nil {
		return nil, lex8.SingleErr(e)
	}

	// codegen
	lib, errs := ir.BuildPkg(b.p)
	if errs != nil {
		return nil, errs
	}

	ret := &build8.Package{
		Lang: "g8-barefunc",
		Main: startName,
		Lib:  lib,
	}
	return ret, nil
}

// CompileBareFunc compiles a bare function into a bare-metal E8 image
func CompileBareFunc(fname, s string) ([]byte, []*lex8.Error, []byte) {
	lang := BareFunc()
	return buildSingle(fname, s, lang)
}
