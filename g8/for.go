package g8

import (
	"e8vm.io/e8vm/g8/tast"
)

func buildForStmt(b *builder, stmt *tast.ForStmt) {
	if stmt.Init != nil {
		b.buildStmt2(stmt.Init)
	}

	if stmt.Cond == nil {
		body := b.f.NewBlock(b.b)
		after := b.f.NewBlock(body)
		body.Jump(body)

		b.b = body
		b.breaks.push(after, "")
		b.continues.push(body, "")

		b.buildStmt2(stmt.Body)

		b.breaks.pop()
		b.continues.pop()

		if stmt.Iter != nil {
			b.buildStmt2(stmt.Iter)
		}
		b.b = after
		return
	}

	condBlock := b.f.NewBlock(b.b)
	body := b.f.NewBlock(condBlock)
	after := b.f.NewBlock(body)
	body.Jump(condBlock)

	b.b = condBlock
	c := b.buildExpr(stmt.Cond)
	b.b.JumpIfNot(c.IR(), after)

	b.b = body
	b.breaks.push(after, "")
	b.continues.push(condBlock, "")

	b.buildStmt2(stmt.Body)

	b.breaks.pop()
	b.continues.pop()

	if stmt.Iter != nil {
		b.buildStmt2(stmt.Iter)
	}

	b.b = after
}
