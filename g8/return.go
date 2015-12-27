package g8

import (
	"e8vm.io/e8vm/g8/tast"
)

func buildReturnStmt(b *builder, stmt *tast.ReturnStmt) {
	if stmt.Exprs != nil {
		exprs := b.buildExpr2(stmt.Exprs)
		assign(b, b.fretRef, exprs)
	}

	next := b.f.NewBlock(b.b)
	b.b.Jump(b.f.End())
	b.b = next
}
