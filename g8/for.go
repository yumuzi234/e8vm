package g8

import (
	"e8vm.io/e8vm/g8/tast"
)

func buildForStmt(b *builder, stmt *tast.ForStmt) {
	if stmt.Init != nil {
		b.buildStmt(stmt.Init)
	}

	if stmt.Cond == nil {
		body := b.f.NewBlock(b.b)
		iter := b.f.NewBlock(body)
		after := b.f.NewBlock(iter)
		body.Jump(body)

		b.b = body
		b.breaks.push(after, "")
		b.continues.push(iter, "")

		b.buildStmt(stmt.Body)

		b.breaks.pop()
		b.continues.pop()

		b.b = iter
		if stmt.Iter != nil {
			b.buildStmt(stmt.Iter)
		}
		b.b = after
		return
	}

	condBlock := b.f.NewBlock(b.b)
	body := b.f.NewBlock(condBlock)
	iter := b.f.NewBlock(body)
	after := b.f.NewBlock(iter)
	iter.Jump(condBlock)

	b.b = condBlock
	c := b.buildExpr(stmt.Cond)
	b.b.JumpIfNot(c.IR(), after)

	b.b = body
	b.breaks.push(after, "")
	b.continues.push(iter, "")

	b.buildStmt(stmt.Body)

	b.breaks.pop()
	b.continues.pop()

	b.b = iter
	if stmt.Iter != nil {
		b.buildStmt(stmt.Iter)
	}

	b.b = after
}

func buildContinueStmt(b *builder) {
	after := b.f.NewBlock(b.b)
	b.b.Jump(b.continues.top())
	b.b = after
}

func buildBreakStmt(b *builder) {
	after := b.f.NewBlock(b.b)
	b.b.Jump(b.breaks.top())
	b.b = after
}
