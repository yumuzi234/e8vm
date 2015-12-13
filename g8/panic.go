package g8

import (
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/g8/ir"
)

func callPanic(b *builder, msg string) {
	if b.panicFunc == nil {
		panic("panic function missing")
	}
	// TODO: print a message
	b.b.Call(nil, b.panicFunc, types.VoidFunc.Sig)
}
