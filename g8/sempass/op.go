package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildConstOpExpr(b *Builder, expr *ast.OpExpr) tast.Expr {
	if expr.A == nil {
		return buildConstUnaryOpExpr(b, expr)
	}
	return buildConstBinaryOpExpr(b, expr)
}

func buildConstUnaryOpExpr(b *Builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op

	B := b.BuildConstExpr(expr.B)
	if B == nil {
		return nil
	}
	return unaryOpConst(b, opTok, B)
}

func buildConstBinaryOpExpr(b *Builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op

	A := b.BuildConstExpr(expr.A)
	if A == nil {
		return nil
	}

	B := b.BuildConstExpr(expr.B)
	if B == nil {
		return nil
	}
	return binaryOpConst(b, opTok, A, B)
}
