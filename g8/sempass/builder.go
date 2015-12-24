package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

type builder struct {
	*lex8.ErrorList
	path string

	this  *tast.Ref
	scope *sym8.Scope

	exprFunc func(b *builder, expr ast.Expr) tast.Expr
	stmtFunc func(b *builder, stmt ast.Stmt) tast.Stmt
}

func newBuilder(path string) *builder {
	return &builder{
		ErrorList: lex8.NewErrorList(),
		path:      path,
		scope:     sym8.NewScope(),
	}
}

func (b *builder) buildExpr(expr ast.Expr) tast.Expr {
	return b.exprFunc(b, expr)
}

func (b *builder) BuildExpr(expr ast.Expr) tast.Expr {
	return b.buildExpr(expr)
}
