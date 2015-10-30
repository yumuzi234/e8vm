package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

func buildForStmt(b *builder, stmt *ast.ForStmt) {
	b.scope.Push()
	defer b.scope.Pop()

	if stmt.Init != nil {
		b.buildStmt(stmt.Init)
	}

	if stmt.Cond == nil {
		body := b.f.NewBlock(b.b)
		after := b.f.NewBlock(body)
		body.Jump(body)

		b.b = body

		b.breaks.push(after, "")
		b.continues.push(body, "")

		buildBlock(b, stmt.Body)

		b.breaks.pop()
		b.continues.pop()

		if stmt.Iter != nil {
			b.buildStmt(stmt.Iter)
		}

		b.b = after
	} else {
		condBlock := b.f.NewBlock(b.b)
		body := b.f.NewBlock(condBlock)
		after := b.f.NewBlock(body)
		body.Jump(condBlock)

		b.b = condBlock
		c := b.buildExpr(stmt.Cond)
		if c == nil {
			return
		}

		if !c.IsBool() {
			pos := ast.ExprPos(stmt.Cond)
			b.Errorf(pos, "expect boolean expression, got %s", c)
			b.b = after
			return
		}
		b.b.JumpIfNot(c.IR(), after)

		b.b = body

		b.breaks.push(after, "")
		b.continues.push(condBlock, "")

		buildBlock(b, stmt.Body)

		b.breaks.pop()
		b.continues.pop()

		if stmt.Iter != nil {
			b.buildStmt(stmt.Iter)
		}

		b.b = after
	}
}
