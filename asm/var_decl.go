package asm

import (
	"shanhu.io/smlvm/asm/ast"
	"shanhu.io/smlvm/lexing"
)

type varDecl struct {
	*ast.Var

	stmts []*varStmt
}

func resolveVar(log lexing.Logger, v *ast.Var) *varDecl {
	ret := new(varDecl)

	ret.Var = v

	for _, stmt := range v.Stmts {
		r := resolveVarStmt(log, stmt)
		ret.stmts = append(ret.stmts, r)
	}

	return ret
}
