package g8

import (
	"e8vm.io/e8vm/g8/types"
)

func buildShift(b *builder, ret, A, B *ref, op string) {
	if types.IsSigned(A.Type()) {
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
	} else {
		if op == ">>" {
			b.b.Arith(ret.IR(), A.IR(), "u>>", B.IR())
		} else {
			b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		}
	}
}
