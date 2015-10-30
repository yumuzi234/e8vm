package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func binaryOpNil(b *builder, opTok *lex8.Token, A, B *ref) *ref {
	op := opTok.Lit
	switch op {
	case "==":
		// nil == nil returns true
		return refTrue
	case "!=":
		return refFalse
	}

	b.Errorf(opTok.Pos, "%q on pointers", op)
	return nil
}

func binaryOpPtr(b *builder, opTok *lex8.Token, A, B *ref) *ref {
	op := opTok.Lit
	atyp := A.Type()
	btyp := B.Type()

	// replace nil with a typed zero
	switch op {
	case "==", "!=":
		if types.IsNil(atyp) {
			A = newRef(btyp, ir.Num(0))
		} else if types.IsNil(btyp) {
			B = newRef(atyp, ir.Num(0))
		}

		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	}

	b.Errorf(opTok.Pos, "%q on pointers", op)
	return nil
}

func testNilSlice(b *builder, r *ref, neg bool) *ref {
	addr := b.newPtr()
	isNil := b.newCond()
	b.b.Arith(addr, nil, "&", r.IR())
	b.b.Arith(isNil, nil, "?", ir.NewAddrRef(addr, 4, 0, false, true))
	if neg {
		b.b.Arith(isNil, nil, "!", isNil)
	}
	return newRef(types.Bool, isNil)
}

func binaryOpSlice(b *builder, opTok *lex8.Token, A, B *ref) *ref {
	op := opTok.Lit
	atyp := A.Type()
	btyp := B.Type()

	switch op {
	case "==", "!=":
		if types.IsNil(atyp) {
			return testNilSlice(b, B, op == "==")
		} else if types.IsNil(btyp) {
			return testNilSlice(b, A, op == "==")
		}

		addrA := b.newPtr()
		addrB := b.newPtr()
		b.b.Arith(addrA, nil, "&", A.IR())
		b.b.Arith(addrB, nil, "&", B.IR())
		baseA := ir.NewAddrRef(addrA, 4, 0, false, true)
		sizeA := ir.NewAddrRef(addrA, 4, 4, false, true)
		baseB := ir.NewAddrRef(addrB, 4, 0, false, true)
		sizeB := ir.NewAddrRef(addrB, 4, 4, false, true)

		ptrEq := b.newCond()
		sizeEq := b.newCond()

		b.b.Arith(ptrEq, baseA, "==", baseB)
		b.b.Arith(sizeEq, sizeA, "==", sizeB)

		ret := b.newCond()
		b.b.Arith(ret, ptrEq, "&", sizeEq)
		if op == "!=" {
			b.b.Arith(ret, nil, "!", ret)
		}
		return newRef(types.Bool, ret)
	}

	b.Errorf(opTok.Pos, "%q on slices", op)
	return nil
}

func refAddress(b *builder, opTok *lex8.Token, B *ref) *ref {
	ret := b.newTemp(&types.Pointer{B.Type()})
	op := opTok.Lit
	opPos := opTok.Pos

	if B.IsType() || !B.IsSingle() {
		b.Errorf(opPos, "invalid operation: %q on %s", op, B)
		return nil
	} else if !B.Addressable() {
		b.Errorf(opPos, "reading address of non-addressable")
		return nil
	}
	b.b.Arith(ret.IR(), nil, op, B.IR())
	return ret
}
