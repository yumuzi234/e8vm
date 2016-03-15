package g8

import (
	"e8vm.io/e8vm/g8/codegen"
	"e8vm.io/e8vm/g8/types"
)

func buildBasicArith(b *builder, ret, A, B *ref, op string) {
	switch op {
	case "/", "%", "u/", "u%":
		isZero := b.newCond()
		b.b.Arith(isZero, B.IR(), "==", codegen.Num(0))

		zeroPanic := b.f.NewBlock(b.b)
		after := b.f.NewBlock(zeroPanic)

		b.b.JumpIfNot(isZero, after)
		b.b = zeroPanic
		callPanic(b, "divided by zero")
		b.b = after
	}

	b.b.Arith(ret.IR(), A.IR(), op, B.IR())
}

func binaryOpInt(b *builder, op string, A, B *ref, t types.T) *ref {
	switch op {
	case "+", "-", "*", "&", "|", "^", "%", "/":
		ret := b.newTemp(t)
		buildBasicArith(b, ret, A, B, op)
		return ret
	case "==", "!=", ">", "<", ">=", "<=":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	}
	panic("bug")
}

func binaryOpUint(b *builder, op string, A, B *ref, t types.T) *ref {
	switch op {
	case "+", "-", "&", "|", "^":
		ret := b.newTemp(t)
		buildBasicArith(b, ret, A, B, op)
		return ret
	case "*", "%", "/":
		ret := b.newTemp(t)
		buildBasicArith(b, ret, A, B, "u"+op)
		return ret
	case "==", "!=":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	case ">", "<", ">=", "<=":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), "u"+op, B.IR())
		return ret
	}
	panic("bug")
}

func binaryOpConst(b *builder, op string, A, B *ref) *ref {
	va, _ := types.NumConst(A.Type())
	vb, _ := types.NumConst(B.Type())

	br := func(b bool) *ref {
		if b {
			return refTrue
		}
		return refFalse
	}
	switch op {
	case "==":
		return br(va == vb)
	case "!=":
		return br(va != vb)
	case ">":
		return br(va > vb)
	case "<":
		return br(va < vb)
	case ">=":
		return br(va >= vb)
	case "<=":
		return br(va <= vb)
	}
	panic("bug")
}

func unaryOpInt(b *builder, op string, B *ref) *ref {
	switch op {
	case "+":
		return B
	case "-", "^":
		ret := b.newTemp(B.Type())
		b.b.Arith(ret.IR(), nil, op, B.IR())
		return ret
	}
	panic("bug")
}
