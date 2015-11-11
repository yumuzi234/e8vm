package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func constNumIr(v int64, t types.T) ir.Ref {
	b, ok := t.(types.Basic)
	if ok {
		switch b {
		case types.Int:
			return ir.Snum(int32(v))
		case types.Uint:
			return ir.Num(uint32(v))
		case types.Int8, types.Uint8:
			println(v)
			return ir.Byt(uint8(v))
		}
	}
	panic("expect an integer type")
}

func constCast(
	b *builder, pos *lex8.Pos, v int64, to types.T,
) *ref {
	if types.IsInteger(to) && types.InRange(v, to) {
		return newRef(to, constNumIr(v, to))
	}

	b.Errorf(pos, "cannot cast %d to %s", v, to)
	return nil
}

func constCastInt(b *builder, pos *lex8.Pos, v int64) *ref {
	return constCast(b, pos, v, types.Int)
}

func constCastUint(b *builder, pos *lex8.Pos, v int64) *ref {
	return constCast(b, pos, v, types.Uint)
}
