package asm

import (
	"io"

	"shanhu.io/smlvm/asm/parse"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
)

func listImport(
	f string, rc io.ReadCloser, imp builds.Importer,
) []*lexing.Error {
	astFile, es := parse.File(f, rc)
	if es != nil {
		return es
	}

	if astFile.Imports == nil {
		return nil
	}

	log := lexing.NewErrorList()
	impDecl := resolveImportDecl(log, astFile.Imports)
	if es := log.Errs(); es != nil {
		return es
	}

	for as, stmt := range impDecl.stmts {
		imp.Import(as, stmt.path, stmt.Path.Pos)
	}

	return nil
}
