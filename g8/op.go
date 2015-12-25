package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
)

func buildBinaryOpExpr(b *builder, expr *ast.OpExpr) *ref {
	opTok := expr.Op
	op := opTok.Lit
	opPos := opTok.Pos

	A := b.buildExpr(expr.A)
	if A == nil {
		return nil
	} else if !A.IsSingle() {
		b.Errorf(opPos, "%q on expression list", op)
		return nil
	}
	atyp := A.Type()

	if types.IsBasic(atyp, types.Bool) && (op == "&&" || op == "||") {
		switch op {
		case "&&":
			return buildLogicAnd(b, opTok, A, expr.B)
		case "||":
			return buildLogicOr(b, opTok, A, expr.B)
		}
		panic("unreachable")
	}

	B := b.buildExpr(expr.B)
	if B == nil { // some error occured
		return nil
	} else if !B.IsSingle() {
		b.Errorf(opPos, "%q on expression list", op)
		return nil
	}
	btyp := B.Type()

	if types.IsConst(atyp) && types.IsConst(btyp) {
		return binaryOpConst(b, opTok, A, B)
	}

	if op == ">>" || op == "<<" {
		if v, ok := types.NumConst(btyp); ok {
			B = constCast(b, opPos, v, types.Uint)
			if B == nil {
				return nil
			}
			btyp = types.Uint
		} else {
			// TODO: typed const
		}
		if v, ok := types.NumConst(atyp); ok {
			A = constCast(b, opPos, v, types.Int)
			if A == nil {
				return nil
			}
			atyp = types.Int
		} else {
			// TODO: typed const
		}
		if !canShift(b, atyp, btyp, opPos, op) {
			return nil
		}

		ret := b.newTemp(atyp)
		buildShift(b, ret, A, B, op)
		return ret
	}

	if v, ok := types.NumConst(atyp); ok {
		A = constCast(b, opPos, v, btyp)
		if A == nil {
			return nil
		}
		atyp = btyp
	} else if c, ok := atyp.(*types.Const); ok {
		atyp = c.Type
	}

	if v, ok := types.NumConst(btyp); ok {
		B = constCast(b, opPos, v, atyp)
		if B == nil {
			return nil
		}
		btyp = atyp
	} else if c, ok := btyp.(*types.Const); ok {
		btyp = c.Type
	}

	if ok, t := types.SameBasic(atyp, btyp); ok {
		switch t {
		case types.Int, types.Int8:
			return binaryOpInt(b, opTok, A, B, t)
		case types.Uint, types.Uint8:
			return binaryOpUint(b, opTok, A, B, t)
		case types.Bool:
			return binaryOpBool(b, opTok, A, B)
		case types.Float32:
			b.Errorf(opPos, "floating point operations not implemented")
			return nil
		}
	}

	// TODO: the branches here are fundamentally upcasting nil pointer to
	// pointer, slice or func pointer this should be done in a better way that
	// is similar to upcasting consts.
	if types.IsNil(atyp) && types.IsNil(btyp) {
		return binaryOpNil(b, opTok, A, B)
	} else if types.BothPointer(atyp, btyp) {
		return binaryOpPtr(b, opTok, A, B)
	} else if types.BothFuncPointer(atyp, btyp) {
		return binaryOpPtr(b, opTok, A, B)
	} else if types.BothSlice(atyp, btyp) {
		return binaryOpSlice(b, opTok, A, B)
	}

	b.Errorf(opPos, "invalid %q", op)
	return nil
}

func buildUnaryOpExpr(b *builder, expr *ast.OpExpr) *ref {
	opTok := expr.Op
	op := opTok.Lit
	opPos := opTok.Pos

	B := b.buildExpr(expr.B)
	if B == nil {
		return nil
	} else if !B.IsSingle() {
		b.Errorf(opPos, "%q on expression list", op)
		return nil
	}

	btyp := B.Type()
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

func buildOpExpr(b *builder, expr *ast.OpExpr) *ref {
	if expr.A == nil {
		return buildUnaryOpExpr(b, expr)
	}
	return buildBinaryOpExpr(b, expr)
}
