package pl

import (
	"e8vm.io/e8vm/pl/codegen"
	"e8vm.io/e8vm/pl/tast"
	"e8vm.io/e8vm/pl/types"
)

func buildStarExpr(b *builder, expr *tast.StarExpr) *ref {
	addr := b.buildExpr(expr.Expr)
	nilPointerPanic(b, addr.IR())
	t := addr.Type().(*types.Pointer).T
	retIR := codegen.NewAddrRef(
		addr.IR(), // base
		t.Size(),  // size
		0,         // offset
		types.IsBasic(t, types.Uint8), // is byte?
		t.RegSizeAlign(),              // is aligned?
	)
	return newAddressableRef(t, retIR)
}
