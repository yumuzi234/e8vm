package asm

import (
	"io"

	"shanhu.io/smlvm/asm/parse"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
)

func listImport(
	f string, rc io.ReadCloser, lst *builds.ImportList,
) []*lexing.Error {
	astFile, errs := parse.File(f, rc)
	if errs != nil {
		return errs
	}

	if astFile.Imports == nil {
		return nil
	}

	log := lexing.NewErrorList()
	impDecl := resolveImportDecl(log, astFile.Imports)
	if errs := log.Errs(); errs != nil {
		return errs
	}

	for as, stmt := range impDecl.stmts {
		lst.Add(as, stmt.path, stmt.Path.Pos)
	}

	return nil
}
