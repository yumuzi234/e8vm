package sempass

import (
	"e8vm.io/e8vm/glang/types"
	"e8vm.io/e8vm/lexing"
)

func canShift(b *builder, atyp, btyp types.T, pos *lexing.Pos, op string) bool {
	if !types.IsInteger(atyp) {
		b.Errorf(pos, "%q on %s", op, atyp)
		return false
	} else if !types.IsInteger(btyp) {
		b.Errorf(pos, "%q with %s", op, btyp)
		return false
	} else if !types.IsUnsigned(btyp) {
		b.Errorf(pos, "%q with %s; must be unsigned", op, btyp)
		return false
	}
	return true
}
