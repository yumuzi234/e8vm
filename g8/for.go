package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func genForStmt(b *builder, stmt *tast.ForStmt) {
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
	c := b.buildExpr2(stmt.Cond)
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

func buildForStmt(b *builder, stmt *ast.ForStmt) {
	b.scope.Push()
	defer b.scope.Pop()

	if stmt.Init != nil {
		b.buildStmt(stmt.Init)
	}

	if stmt.Cond == nil { // infinite for
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
		return
	}

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
