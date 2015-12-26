package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
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
			return ir.Byt(uint8(v))
		}
	}
	panic("expect an integer type")
}
