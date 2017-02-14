package pl

import (
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildSwitchStmt(b *builder, stmt *tast.SwitchStmt) {
	type caseInfo struct {
		expr *codegen.Block
		body *codegen.Block
		c    *tast.Case
	}

	s := buildExpr(b, stmt.Expr)
	cases := make([]*caseInfo, len(stmt.Cases))

	start := b.b
	def := b.f.NewBlock(start)
	after := b.f.NewBlock(def)
	for i := len(cases) - 1; i >= 0; i-- {
		c := &caseInfo{c: stmt.Cases[i]}
		if c.c.Expr != nil {
			c.expr = b.f.NewBlock(start)
		} else {
			c.expr = def
		}
		c.body = b.f.NewBlock(def)
		cases[i] = c
	}
	def.Jump(after)

	for _, c := range cases {
		b.b = c.expr
		if c.c.Expr != nil {
			ret := b.newTemp(types.Bool)
			b.b.Arith(ret.IR(), s.IR(), "==", b.buildExpr(c.c.Expr).IR())
			b.b.JumpIf(ret.IR(), c.body)
		} else {
			b.b.Jump(c.body)
		}

		b.b = c.body
		for _, s := range c.c.Stmts {
			b.buildStmt(s)
		}
		if !c.c.Fallthrough {
			b.b.Jump(after)
		}
	}
	b.b = after
}
