package asm

import (
	"shanhu.io/smlvm/asm/ast"
	"shanhu.io/smlvm/lexing"
)

type importDecl struct {
	*ast.Import

	stmts map[string]*importStmt
}

func resolveImportDecl(log lexing.Logger, imp *ast.Import) *importDecl {
	ret := new(importDecl)
	ret.Import = imp
	ret.stmts = make(map[string]*importStmt)

	for _, stmt := range imp.Stmts {
		r := resolveImportStmt(log, stmt)

		if other, found := ret.stmts[r.as]; found {
			log.Errorf(r.pos(), "%s already imported", r.as)
			log.Errorf(other.pos(), "  previously imported here")
			continue
		}

		ret.stmts[r.as] = r
	}

	return ret
}
