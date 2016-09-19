package asm

import (
	"shanhu.io/smlvm/asm/ast"
	"shanhu.io/smlvm/lexing"
)

type funcDecl struct {
	*ast.Func

	stmts []*funcStmt
}

func resolveFunc(log lexing.Logger, f *ast.Func) *funcDecl {
	ret := new(funcDecl)
	ret.Func = f

	for _, stmt := range f.Stmts {
		r := resolveFuncStmt(log, stmt)
		ret.stmts = append(ret.stmts, r)
	}

	return ret
}
