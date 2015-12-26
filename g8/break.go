package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

func buildBreakStmt(b *builder, stmt *ast.BreakStmt) {
	if stmt.Label != nil {
		b.Errorf(stmt.Label.Pos, "break with label not implemented")
		return
	}

	next := b.breaks.top()
	if next == nil {
		b.Errorf(stmt.Kw.Pos, "break is not in a for or switch block")
		return
	}
	after := b.f.NewBlock(b.b)
	b.b.Jump(next)
	b.b = after
}

func genBreakStmt(b *builder) {
	after := b.f.NewBlock(b.b)
	b.b.Jump(b.breaks.top())
	b.b = after
}
