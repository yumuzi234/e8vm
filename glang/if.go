package glang

import (
	"e8vm.io/e8vm/glang/tast"
)

func buildIfStmt(b *builder, stmt *tast.IfStmt) {
	c := b.buildExpr(stmt.Expr)
	if stmt.Else == nil {
		body := b.f.NewBlock(b.b)
		after := b.f.NewBlock(body)
		b.b.JumpIfNot(c.IR(), after)
		b.b = body
		b.buildStmt(stmt.Body)
		b.b = after
		return
	}

	ifBody := b.f.NewBlock(b.b)
	elseBody := b.f.NewBlock(ifBody)
	after := b.f.NewBlock(elseBody)
	b.b.JumpIfNot(c.IR(), elseBody)
	ifBody.Jump(after)

	b.b = ifBody // switch to if body
	b.buildStmt(stmt.Body)
	b.b = elseBody
	b.buildStmt(stmt.Else)
	b.b = after
}
