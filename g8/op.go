package g8

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildOpExpr(b *builder, expr *tast.OpExpr) *ref {
	if expr.A == nil {
		return buildUnaryOpExpr(b, expr)
	}
	return buildBinaryOpExpr(b, expr)
}

func buildUnaryOpExpr(b *builder, expr *tast.OpExpr) *ref {
	op := expr.Op.Lit
	B := b.buildExpr2(expr.B)
	btyp := B.Type()
	if op == "&" {
		ret := b.newTemp(&types.Pointer{btyp})
		b.b.Arith(ret.IR(), nil, op, B.IR())
		return ret
	} else if types.IsConst(btyp) {
		panic("bug")
	} else if types.IsInteger(btyp) {
		return unaryOpInt(b, op, B)
	} else if types.IsBasic(btyp, types.Bool) {
		return unaryOpBool(b, op, B)
	}
	panic("bug")
}

func buildBinaryOpExpr(b *builder, expr *tast.OpExpr) *ref {
	op := expr.Op.Lit
	A := b.buildExpr2(expr.A)
	atyp := A.Type()
	if types.IsBasic(atyp, types.Bool) && (op == "&&" || op == "||") {
		switch op {
		case "&&":
			return buildLogicAnd(b, A, expr.B)
		case "||":
			return buildLogicOr(b, A, expr.B)
		}
		panic("unreachable")
	}

	B := b.buildExpr2(expr.B)
	btyp := B.Type()
	if types.IsConst(atyp) && types.IsConst(btyp) {
		return binaryOpConst(b, op, A, B)
	}

	if op == ">>" || op == "<<" {
		ret := b.newTemp(atyp)
		buildShift(b, ret, A, B, op)
		return ret
	}

	if ok, t := types.SameBasic(atyp, btyp); ok {
		switch t {
		case types.Int, types.Int8:
			return binaryOpInt(b, op, A, B, t)
		case types.Uint, types.Uint8:
			return binaryOpUint2(b, op, A, B, t)
		case types.Bool:
			return binaryOpBool(b, op, A, B)
		}
		panic("bug")
	}

	if types.IsNil(atyp) && types.IsNil(btyp) {
		return binaryOpNil(b, op, A, B)
	} else if types.BothPointer(atyp, btyp) {
		return binaryOpPtr(b, op, A, B)
	} else if types.BothFuncPointer(atyp, btyp) {
		return binaryOpPtr(b, op, A, B)
	} else if types.BothSlice(atyp, btyp) {
		return binaryOpSlice(b, op, A, B)
	}
	panic("bug")
}
