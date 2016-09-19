package pl

import (
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/types"
)

func binaryOpNil(b *builder, op string, A, B *ref) *ref {
	switch op {
	case "==":
		return refTrue
	case "!=":
		return refFalse
	}
	panic("bug")
}

func binaryOpPtr(b *builder, op string, A, B *ref) *ref {
	atyp := A.Type()
	btyp := B.Type()

	switch op {
	case "==", "!=":
		// replace nil with a typed zero
		if types.IsNil(atyp) {
			A = newRef(btyp, codegen.Num(0))
		} else if types.IsNil(btyp) {
			B = newRef(atyp, codegen.Num(0))
		}

		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	}
	panic("bug")
}

func testNilSlice(b *builder, r *ref, neg bool) *ref {
	addr := b.newPtr()
	isNil := b.newCond()
	b.b.Arith(addr, nil, "&", r.IR())
	b.b.Arith(isNil, nil, "?", codegen.NewAddrRef(addr, 4, 0, false, true))
	if neg {
		b.b.Arith(isNil, nil, "!", isNil)
	}
	return newRef(types.Bool, isNil)
}

func binaryOpSlice(b *builder, op string, A, B *ref) *ref {
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
		baseA := codegen.NewAddrRef(addrA, 4, 0, false, true)
		sizeA := codegen.NewAddrRef(addrA, 4, 4, false, true)
		baseB := codegen.NewAddrRef(addrB, 4, 0, false, true)
		sizeB := codegen.NewAddrRef(addrB, 4, 4, false, true)

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
	panic("bug")
}
