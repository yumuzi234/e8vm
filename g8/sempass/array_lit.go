package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildArrayLit(b *builder, lit *ast.ArrayLiteral) tast.Expr {
	b.Errorf(ast.ExprPos(lit), "array literal not implemented")
	return nil
}
