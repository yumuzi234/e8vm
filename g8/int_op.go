package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func buildBasicArith(b *builder, ret, A, B *ref, op string) {
	if op == "%" || op == "/" {
		isZero := b.newCond()
		b.b.Arith(isZero, B.IR(), "==", ir.Num(0))

		zeroPanic := b.f.NewBlock(b.b)
		after := b.f.NewBlock(zeroPanic)

		b.b.JumpIfNot(isZero, after)
		b.b = zeroPanic
		callPanic(b, "divided by zero")
		b.b = after
	}

	b.b.Arith(ret.IR(), A.IR(), op, B.IR())
}

func binaryOpInt(b *builder, opTok *lex8.Token, A, B *ref, t types.T) *ref {
	op := opTok.Lit
	switch op {
	case "+", "-", "*", "&", "|", "^", "%", "/":
		ret := b.newTemp(t)
		buildBasicArith(b, ret, A, B, op)
		return ret
		return ret
	case "==", "!=", ">", "<", ">=", "<=":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	}

	b.Errorf(opTok.Pos, "%q on ints", op)
	return nil
}

func binaryOpUint(b *builder, opTok *lex8.Token, A, B *ref, t types.T) *ref {
	op := opTok.Lit
	switch op {
	case "+", "-", "*", "&", "|", "^", "%", "/":
		ret := b.newTemp(t)
		buildBasicArith(b, ret, A, B, op)
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

	b.Errorf(opTok.Pos, "%q on ints", op)
	return nil
}

func binaryOpConst(b *builder, opTok *lex8.Token, A, B *ref) *ref {
	op := opTok.Lit
	if !A.IsSingle() || !B.IsSingle() {
		b.Errorf(opTok.Pos, "invalid %s %q %s", A, op, B)
		return nil
	}

	va, oka := types.NumConst(A.Type())
	vb, okb := types.NumConst(B.Type())
	if !(oka && okb) {
		b.Errorf(opTok.Pos, "non-numeric consts ops not implemented")
		return nil
	}

	r := func(v int64) *ref {
		return newRef(types.NewNumber(v), nil)
	}
	br := func(b bool) *ref {
		if b {
			return refTrue
		}
		return refFalse
	}

	switch op {
	case "+":
		return r(va + vb)
	case "-":
		return r(va - vb)
	case "*":
		return r(va * vb)
	case "&":
		return r(va & vb)
	case "|":
		return r(va | vb)
	case "^":
		return r(va ^ vb)
	case "%":
		if vb == 0 {
			b.Errorf(opTok.Pos, "modular by zero")
			return nil
		}
		return r(va % vb)
	case "/":
		if vb == 0 {
			b.Errorf(opTok.Pos, "divide by zero")
			return nil
		}
		return r(va / vb)
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
	case "<<":
		if vb < 0 {
			b.Errorf(opTok.Pos, "shift with negative value", vb)
			return nil
		}
		return r(va << uint64(vb))
	case ">>":
		if vb < 0 {
			b.Errorf(opTok.Pos, "shift with negative value", vb)
			return nil
		}
		return r(va >> uint64(vb))
	}

	b.Errorf(opTok.Pos, "%q on consts", op)
	return nil
}

func unaryOpInt(b *builder, opTok *lex8.Token, B *ref) *ref {
	op := opTok.Lit
	switch op {
	case "+":
		return B
	case "-", "^":
		ret := b.newTemp(B.Type())
		b.b.Arith(ret.IR(), nil, op, B.IR())
		return ret
	}

	b.Errorf(opTok.Pos, "invalid operation: %q on %s", op, B)
	return nil
}

func unaryOpConst(b *builder, opTok *lex8.Token, B *ref) *ref {
	op := opTok.Lit
	if !B.IsSingle() {
		b.Errorf(opTok.Pos, "invalid operation: %q on %s", op, B)
		return nil
	}

	v, ok := types.NumConst(B.Type())
	if !ok {
		// TODO
		b.Errorf(opTok.Pos, "typed operation not implemented")
		return nil
	}

	switch op {
	case "+":
		return B
	case "-":
		return newRef(types.NewNumber(-v), nil)
	}

	b.Errorf(opTok.Pos, "invalid operation: %q on %s", op, B)
	return nil
}
