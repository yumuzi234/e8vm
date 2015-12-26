package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

func buildContinueStmt(b *builder, stmt *ast.ContinueStmt) {
	if stmt.Label != nil {
		b.Errorf(stmt.Label.Pos, "continue with label not implemented")
		return
	}

	next := b.continues.top()
	if next == nil {
		b.Errorf(stmt.Kw.Pos, "continue is not in a for block")
		return
	}

	after := b.f.NewBlock(b.b)
	b.b.Jump(next)
	b.b = after
}

func genContinueStmt(b *builder) {
	after := b.f.NewBlock(b.b)
	b.b.Jump(b.continues.top())
	b.b = after
}
