package pl

import (
	"e8vm.io/e8vm/pl/tast"
)

func buildReturnStmt(b *builder, stmt *tast.ReturnStmt) {
	if stmt.Exprs != nil {
		exprs := b.buildExpr(stmt.Exprs)
		assign(b, b.fretRef, exprs)
	}

	next := b.f.NewBlock(b.b)
	b.b.Jump(b.f.End())
	b.b = next
}
