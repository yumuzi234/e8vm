package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func etSize(t types.T) int32 {
	ret := t.Size()
	if t.RegSizeAlign() {
		return types.RegSizeAlignUp(ret)
	}
	return ret
}

func checkInRange(b *builder, index, n ir.Ref, op string) {
	inRange := b.newCond()
	b.b.Arith(inRange, index, op, n)

	outOfRange := b.f.NewBlock(b.b)
	after := b.f.NewBlock(outOfRange)
	b.b.JumpIf(inRange, after)

	b.b = outOfRange
	callPanic(b, "index out of range")

	b.b = after
}

func newSlice(b *builder, t types.T, addr, size ir.Ref) *ref {
	ret := b.newTemp(&types.Slice{T: t})
	retAddr := b.newPtr()
	b.b.Arith(retAddr, nil, "&", ret.IR())
	b.b.Assign(ir.NewAddrRef(retAddr, 4, 0, false, true), addr)
	b.b.Assign(ir.NewAddrRef(retAddr, 4, 4, false, true), size)
	return ret
}

func genIndexExpr(b *builder, expr *tast.IndexExpr) *ref {
	array := b.buildExpr2(expr.Array)

	if expr.HasColon {
		return genSlicing(b, expr, array)
	}
	return genArrayGet(b, expr, array)
}

func loadArray2(b *builder, array *ref) (addr, n ir.Ref, et types.T) {
	base := b.newPtr()
	t := array.Type()
	switch t := t.(type) {
	case *types.Array:
		b.b.Arith(base, nil, "&", array.IR())
		return base, ir.Snum(t.N), t.T
	case *types.Slice:
		b.b.Arith(base, nil, "&", array.IR())
		addr = ir.NewAddrRef(base, 4, 0, false, true)
		n = ir.NewAddrRef(base, 4, 4, false, true)
		return addr, n, t.T
	}
	panic("bug")
}

func checkArrayIndex2(b *builder, index *ref) ir.Ref {
	t := index.Type()
	if types.IsSigned(t) {
		neg := b.newCond()
		b.b.Arith(neg, nil, "<0", index.IR())
		negPanic := b.f.NewBlock(b.b)
		after := b.f.NewBlock(negPanic)
		b.b.JumpIfNot(neg, after)

		b.b = negPanic
		callPanic(b, "index is negative")

		b.b = after
	}
	return index.IR()
}

func genArrayIndex(b *builder, expr tast.Expr) ir.Ref {
	index := b.buildExpr2(expr)
	return checkArrayIndex2(b, index)
}

func genSlicing(b *builder, expr *tast.IndexExpr, array *ref) *ref {
	baseAddr, n, et := loadArray2(b, array)

	var addr, indexStart, offset ir.Ref
	if expr.Index == nil {
		indexStart = ir.Num(0)
		addr = baseAddr
	} else {
		indexStart = genArrayIndex(b, expr.Index)
		checkInRange(b, indexStart, n, "u<=")

		offset = b.newPtr()
		b.b.Arith(offset, indexStart, "*", ir.Snum(etSize(et)))
		addr = b.newPtr()
		b.b.Arith(addr, baseAddr, "+", offset)
	}

	var indexEnd ir.Ref
	if expr.IndexEnd == nil {
		indexEnd = n
	} else {
		indexEnd = genArrayIndex(b, expr.IndexEnd)
		checkInRange(b, indexEnd, n, "u<=")
		checkInRange(b, indexStart, indexEnd, "u<=")
	}

	size := b.newPtr()
	b.b.Arith(size, indexEnd, "-", indexStart)
	return newSlice(b, et, addr, size)
}

func genArrayGet(b *builder, expr *tast.IndexExpr, array *ref) *ref {
	index := genArrayIndex(b, expr.Index)
	base, n, et := loadArray2(b, array)
	checkInRange(b, index, n, "u<")

	addr := b.newPtr()
	b.b.Arith(addr, index, "*", ir.Snum(etSize(et)))
	b.b.Arith(addr, base, "+", addr)
	size := et.Size()

	retIR := ir.NewAddrRef(
		addr,             // base address
		size,             // size
		0,                // dynamic offset; precalculated
		types.IsByte(et), // isByte
		true,             // isAlign
	)
	return newAddressableRef(et, retIR)
}
