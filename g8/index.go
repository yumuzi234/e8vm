package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func etSize(t types.T) int32 {
	ret := t.Size()
	if t.RegSizeAlign() {
		return types.RegSizeAlignUp(ret)
	}
	return ret
}

func loadArray(b *builder, array *ref) (addr, n ir.Ref, et types.T) {
	base := b.newPtr()
	t := array.Type()
	switch t := t.(type) {
	case *types.Array:
		b.b.Arith(base, nil, "&", array.IR())
		return base, ir.Num(uint32(t.N)), t.T
	case *types.Slice:
		b.b.Arith(base, nil, "&", array.IR())
		addr = ir.NewAddrRef(base, 4, 0, false, true)
		n = ir.NewAddrRef(base, 4, 4, false, true)
		return addr, n, t.T
	default:
		return nil, nil, nil
	}
}

func checkArrayIndex(b *builder, index *ref, pos *lex8.Pos) ir.Ref {
	t := index.Type()
	if v, ok := types.NumConst(t); ok {
		if v < 0 {
			b.Errorf(pos, "array index is negative: %d", v)
			return nil
		}

		index = constCastInt(b, pos, v)
		if index == nil {
			return nil
		}
		return index.IR()
	}

	if !types.IsInteger(t) {
		b.Errorf(pos, "index must be an integer")
		return nil
	}

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

func buildArrayIndex(b *builder, expr ast.Expr, pos *lex8.Pos) ir.Ref {
	index := b.buildExpr(expr)
	if index == nil {
		return nil
	} else if !index.IsSingle() {
		b.Errorf(pos, "index with expression list")
		return nil
	}

	return checkArrayIndex(b, index, pos)
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

func buildSlicing(b *builder, expr *ast.IndexExpr, array *ref) *ref {
	baseAddr, n, et := loadArray(b, array)
	if et == nil {
		b.Errorf(expr.Lbrack.Pos, "slicing on neither array nor slice")
		return nil
	}

	var addr, indexStart, offset ir.Ref
	if expr.Index == nil {
		indexStart = ir.Num(0)
		addr = baseAddr
	} else {
		indexStart = buildArrayIndex(b, expr.Index, expr.Lbrack.Pos)
		checkInRange(b, indexStart, n, "u<=")

		offset = b.newPtr()
		b.b.Arith(offset, indexStart, "*", ir.Num(uint32(etSize(et))))
		addr = b.newPtr()
		b.b.Arith(addr, baseAddr, "+", offset)
	}

	var indexEnd ir.Ref
	if expr.IndexEnd == nil {
		indexEnd = n
	} else {
		indexEnd = buildArrayIndex(b, expr.IndexEnd, expr.Colon.Pos)
		checkInRange(b, indexEnd, n, "u<=")
		checkInRange(b, indexStart, indexEnd, "u<=")
	}

	size := b.newPtr()
	b.b.Arith(size, indexEnd, "-", indexStart)

	return newSlice(b, et, addr, size)
}

func buildArrayGet(b *builder, expr *ast.IndexExpr, array *ref) *ref {
	index := buildArrayIndex(b, expr.Index, expr.Lbrack.Pos)
	if index == nil {
		return nil
	}

	base, n, et := loadArray(b, array)
	if et == nil {
		b.Errorf(expr.Lbrack.Pos, "index on neither array or slice")
		return nil
	}

	checkInRange(b, index, n, "u<")

	addr := b.newPtr()

	b.b.Arith(addr, index, "*", ir.Num(uint32(etSize(et))))
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

func buildIndexExpr(b *builder, expr *ast.IndexExpr) *ref {
	array := b.buildExpr(expr.Array)
	if array == nil {
		return nil
	} else if !array.IsSingle() {
		b.Errorf(expr.Lbrack.Pos, "index on expression list")
		return nil
	}

	if expr.Colon != nil {
		return buildSlicing(b, expr, array)
	}

	return buildArrayGet(b, expr, array)
}
