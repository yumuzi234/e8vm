package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
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

func buildOpExpr(b *Builder, expr *ast.OpExpr) tast.Expr {
	if expr.A == nil {
		return buildUnaryOpExpr(b, expr)
	}
	panic("todo")
	// return buildBinaryOpExpr(b, expr)
}

func buildUnaryOpExpr(b *Builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op
	op := opTok.Lit
	opPos := opTok.Pos

	B := b.BuildExpr(expr.B)
	if B == nil {
		return nil
	}
	bref := tast.ExprRef(B)
	if bref.List != nil {
		b.Errorf(opPos, "%q on expression list", bref.T)
		return nil
	}

	btyp := bref.T
	if op == "&" {
		return refAddress(b, opTok, B)
	} else if types.IsConst(btyp) {
		return unaryOpConst(b, opTok, B)
	} else if types.IsInteger(btyp) {
		return unaryOpInt(b, opTok, B)
	} else if types.IsBasic(btyp, types.Bool) {
		return unaryOpBool(b, opTok, B)
	}

	b.Errorf(opPos, "invalid unary operator %q", op)
	return nil
}
