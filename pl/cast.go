package pl

import (
	"e8vm.io/e8vm/arch"
	"e8vm.io/e8vm/pl/codegen"
	"e8vm.io/e8vm/pl/types"
)

func isPointerType(t types.T) bool {
	return types.IsPointer(t) || types.IsFuncPointer(t)
}

func regSizeCastable(to, from types.T) bool {
	if types.IsPointer(to) && types.IsPointer(from) {
		return true
	}
	if isPointerType(to) && types.IsBasic(from, types.Uint) {
		return true
	}
	if types.IsBasic(to, types.Uint) && isPointerType(from) {
		return true
	}
	return false
}

func constNumIr(v int64, t types.T) codegen.Ref {
	b, ok := t.(types.Basic)
	if ok {
		switch b {
		case types.Int:
			return codegen.Snum(int32(v))
		case types.Uint:
			return codegen.Num(uint32(v))
		case types.Int8:
			return codegen.Byt(uint8(v), false)
		case types.Uint8:
			return codegen.Byt(uint8(v), true)
		}
	}
	panic("expect an integer type")
}

func buildCast(b *builder, from *ref, t types.T) *ref {
	srcType := from.Type()
	ret := b.newTemp(t)

	if types.IsNil(srcType) {
		size := t.Size()
		if size == arch.RegSize {
			return newRef(t, codegen.Num(0))
		}
		if _, ok := t.(*types.Slice); !ok {
			panic("bug")
		}
		ret := b.newTemp(t)
		b.b.Zero(ret.IR())
		return ret
	}

	if c, ok := srcType.(*types.Const); ok {
		if v, ok := types.NumConst(srcType); ok && types.IsInteger(t) {
			return newRef(t, constNumIr(v, t))
		}
		// TODO: we do not support typed const right?
		// so why need this line?
		srcType = c.Type // using the underlying type
	}

	if types.IsInteger(t) && types.IsInteger(srcType) {
		b.b.Arith(ret.IR(), nil, "cast", from.IR())
		return ret
	}
	if regSizeCastable(t, srcType) {
		b.b.Arith(ret.IR(), nil, "", from.IR())
		return ret
	}
	panic("bug")
}
