package pl

import (
	"fmt"

	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/parse"
	"shanhu.io/smlvm/pl/sempass"
)

// because bare function also uses builtin functions that comes from the
// building system, we also need to make it a simple language: a
// language with only one (implicit) main function
// In fact, we can simple "inherit" the basic lang
type bareFunc struct{ *lang }

// BareFunc is a language where it only contains an implicit main function.
func BareFunc() *builds.Lang {
	return &builds.Lang{
		Ext:      "g",
		Compiler: bareFunc{new(lang)},
	}
}

func (bareFunc) Prepare(*builds.FileSet) (
	*builds.ImportList, []*lexing.Error,
) {
	ret := builds.NewImportList()
	ret.Add("$", BuiltInPkg, nil)
	return ret, nil
}

func buildBareFunc(b *builder, stmts []ast.Stmt) []*lexing.Error {
	tstmts, errs := sempass.BuildBareFunc(b.scope, stmts)
	if errs != nil {
		return errs
	}

	b.f = b.p.NewFunc(":start", nil, codegen.VoidFuncSig)
	b.fretRef = nil
	b.b = b.f.NewBlock(nil)

	for _, stmt := range tstmts {
		buildStmt(b, stmt)
	}
	return nil
}

func (bare bareFunc) Compile(pinfo *builds.PkgInfo) (
	pkg *builds.Package, es []*lexing.Error,
) {
	// parsing
	theFile := pinfo.Files.OnlyFile()
	if theFile == nil {
		return nil, lexing.SingleErr(fmt.Errorf("too many files"))
	}
	rc, err := theFile.Open()
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
	stmts, es := parse.Stmts(theFile.Path, rc)
	if es != nil {
		return nil, es
	}

	// building
	b := makeBuilder(pinfo)
	if es = b.Errs(); es != nil {
		return nil, es
	}
	if es := buildBareFunc(b, stmts); es != nil {
		return nil, es
	}
	if es := b.Errs(); es != nil {
		return nil, es
	}

	if e := outputIr(pinfo, b); e != nil {
		return nil, lexing.SingleErr(e)
	}

	// codegen
	lib, errs := codegen.BuildPkg(b.p)
	if errs != nil {
		return nil, errs
	}

	ret := &builds.Package{
		Lang: "g8-barefunc",
		Main: startName,
		Lib:  lib,
	}
	return ret, nil
}

// CompileBareFunc compiles a bare function into a bare-metal E8 image
func CompileBareFunc(fname, s string) ([]byte, []*lexing.Error, []byte) {
	lang := BareFunc()
	return buildSingle(fname, s, lang, new(builds.Options))
}
