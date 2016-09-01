package asm

import (
	"e8vm.io/e8vm/asm/ast"
	"e8vm.io/e8vm/lexing"
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
