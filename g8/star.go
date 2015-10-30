package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

func buildStarExpr(b *builder, expr *ast.StarExpr) *ref {
	opPos := expr.Star.Pos

	addr := b.buildExpr(expr.Expr)
	if addr == nil {
		return nil
	} else if addr.IsType() {
		// a pionter type
		t := addr.TypeType()
		return newTypeRef(&types.Pointer{t})
	} else if !addr.IsSingle() {
		b.Errorf(opPos, "* on expression list")
		return nil
	}

	// referencing the value of a pointer
	t, ok := addr.Type().(*types.Pointer)
	if !ok {
		b.Errorf(opPos, "* on non-pointer")
		return nil
	}
	nilPointerPanic(b, addr.IR())

	retIR := ir.NewAddrRef(
		addr.IR(),  // base
		t.T.Size(), // size
		0,          // offset
		types.IsBasic(t.T, types.Uint8), // is byte?
		t.T.RegSizeAlign(),              // is aligned?
	)
	return newAddressableRef(t.T, retIR)
}
