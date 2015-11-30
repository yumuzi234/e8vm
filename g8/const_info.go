package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

type constInfo struct {
	name *lex8.Token
	typ  ast.Expr
	expr ast.Expr

	deps    []string
	queuing bool
	queued  bool
}

func newConstInfo(name *lex8.Token, typ, expr ast.Expr) *constInfo {
	return &constInfo{
		name: name, typ: typ, expr: expr,
		deps: symUseExpr(expr),
	}
}
