package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
)

func buildSwitchStmt(b *builder, stmt *ast.SwitchStmt) tast.Stmt {
	return buildSwtich(b, stmt.Expr, stmt.Cases)
}

func buildFallthroughStmt(b *builder, stmt *ast.FallthroughStmt) tast.Stmt {
	if b.switchlayer == 0 {
		b.CodeErrorf(stmt.Kw.Pos, "pl.fallthrough.notInLoop",
			"fallthrough is not in a switch block")
		return nil
	}
	return &tast.FallthroughStmt{}
}

func buildSwtich(b *builder, expr ast.Expr, Cases []*ast.Case) tast.Stmt {
	e := b.buildExpr(expr)
	if e == nil {
		return nil
	}
	b.switchlayer++
	var cases []*tast.Case
	for _, c := range Cases {
		ret := buildCase(b, c)
		if ret == nil {
			return nil
		}
		cases = append(cases, ret)
	}
	return &tast.SwitchStmt{Expr: e, Cases: cases}
}

func buildCase(b *builder, c *ast.Case) *tast.Case {
	e := b.buildExpr(c.Expr)
	if e == nil {
		return nil
	}
	var stmts []tast.Stmt
	for _, stmt := range c.Stmts {
		s := b.buildStmt(stmt)
		if s != nil {
			stmts = append(stmts, s)
		}
	}
	return &tast.Case{Expr: e, Stmts: stmts}
}
