package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// Builder defines the intermediate interface for building
type Builder interface {
	BuildExpr(expr ast.Expr) tast.Expr
	Errs() []*lex8.Error
}

// NewBuilder creates a new builder with a specific path.
func NewBuilder(path string, scope *sym8.Scope) Builder {
	ret := newBuilder(path)
	ret.exprFunc = buildExpr
	ret.scope = scope // TODO: remove this

	return ret
}
