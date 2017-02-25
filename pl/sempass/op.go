package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildConstOpExpr(b *builder, expr *ast.OpExpr) tast.Expr {
	if expr.A == nil {
		return buildConstUnaryOpExpr(b, expr)
	}
	return buildConstBinaryOpExpr(b, expr)
}

func buildConstUnaryOpExpr(b *builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op

	B := b.buildConstExpr(expr.B)
	if B == nil {
		return nil
	}
	return unaryOpConst(b, opTok, B)
}

func buildConstBinaryOpExpr(b *builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op

	A := b.buildConstExpr(expr.A)
	if A == nil {
		return nil
	}

	B := b.buildConstExpr(expr.B)
	if B == nil {
		return nil
	}
	return binaryOpConst(b, opTok, A, B)
}

func buildOpExpr(b *builder, expr *ast.OpExpr) tast.Expr {
	hold := b.lhsSwap(false)
	defer b.lhsRestore(hold)

	if expr.A == nil {
		return buildUnaryOpExpr(b, expr)
	}
	return buildBinaryOpExpr(b, expr)
}

func buildUnaryOpExpr(b *builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op
	op := opTok.Lit
	opPos := opTok.Pos

	B := b.buildExpr(expr.B)
	if B == nil {
		return nil
	}
	bref := B.R()
	if !bref.IsSingle() {
		b.Errorf(opPos, "%q on %s", op, bref)
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

func buildBinaryOpExpr(b *builder, expr *ast.OpExpr) tast.Expr {
	opTok := expr.Op
	op := opTok.Lit
	opPos := opTok.Pos

	A := b.buildExpr(expr.A)
	if A == nil {
		return nil
	}
	aref := A.R()
	if !aref.IsSingle() {
		b.Errorf(opPos, "%q on %s", op, aref)
		return nil
	}
	atyp := aref.T

	B := b.buildExpr(expr.B)
	if B == nil {
		return nil
	}
	bref := B.R()
	if !bref.IsSingle() {
		b.Errorf(opPos, "%q on %s", op, bref)
		return nil
	}
	btyp := bref.T

	if types.IsConst(atyp) && types.IsConst(btyp) {
		return binaryOpConst(b, opTok, A, B)
	}

	if op == ">>" || op == "<<" {
		if v, ok := types.NumConst(btyp); ok {
			B = numsCast(b, opPos, v, B, types.Uint)
			if B == nil {
				return nil
			}
			btyp = types.Uint
		}

		if v, ok := types.NumConst(atyp); ok {
			A = numsCast(b, opPos, v, A, types.Int)
			if A == nil {
				return nil
			}
			atyp = types.Int
		}

		if !canShift(b, atyp, btyp, opPos, op) {
			return nil
		}

		r := tast.NewRef(atyp)
		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: r}
	}

	if v, ok := types.NumConst(atyp); ok {
		A = numsCast(b, opPos, v, A, btyp)
		if A == nil {
			return nil
		}
		atyp = btyp
	} else if c, ok := atyp.(*types.Const); ok {
		atyp = c.Type
	}

	if v, ok := types.NumConst(btyp); ok {
		B = numsCast(b, opPos, v, B, atyp)
		if B == nil {
			return nil
		}
		btyp = atyp
	} else if c, ok := btyp.(*types.Const); ok {
		btyp = c.Type
	}

	if ok, t := types.SameBasic(atyp, btyp); ok {
		switch t {
		case types.Int, types.Int8, types.Uint, types.Uint8:
			return binaryOpInt(b, opTok, A, B, t)
		case types.Bool:
			return binaryOpBool(b, opTok, A, B)
		case types.Float32:
			b.Errorf(opPos, "floating point operations not implemented")
			return nil
		}
	}

	if types.IsNil(atyp) && types.IsNil(btyp) {
		return binaryOpNil(b, opTok, A, B)
	} else if types.BothPointer(atyp, btyp) {
		return binaryOpPtr(b, opTok, A, B)
	} else if types.BothFuncPointer(atyp, btyp) {
		return binaryOpPtr(b, opTok, A, B)
	} else if types.BothSlice(atyp, btyp) {
		return binaryOpSlice(b, opTok, A, B)
	}

	b.CodeErrorf(opPos, "pl.invalidOp",
		"invalid operation of %s %s %s", atyp, op, btyp)
	if types.IsInteger(atyp) && types.IsInteger(btyp) {
		switch op {
		case "+", "-", "*", "&", "|", "^", "%", "/",
			"==", "!=", ">", "<", ">=", "<=":
			b.Errorf(
				opPos,
				"operation %s needs the same type on both sides",
				op,
			)
		}
	}
	return nil
}
