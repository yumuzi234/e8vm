package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildSwitchStmt(b *builder, stmt *ast.SwitchStmt) tast.Stmt {
	return buildSwitch(b, stmt.Expr, stmt.Cases)
}

func buildSwitch(b *builder, expr ast.Expr, Cases []*ast.Case) tast.Stmt {
	e := b.buildExpr(expr)
	if e == nil {
		pos := ast.ExprPos(expr)
		b.CodeErrorf(pos, "pl.switchExpr",
			"nil cannot be used as expression for switch")
		return nil
	}
	exprRef := e.R()
	if !exprRef.IsSingle() {
		pos := ast.ExprPos(expr)
		// Q %s
		b.CodeErrorf(pos, "pl.switchExpr",
			"expect single expression for switch, got %s", exprRef)
		return nil
	}
	if !(types.IsInteger(exprRef.Type()) || types.IsConst(exprRef.Type())) {
		pos := ast.ExprPos(expr)
		b.CodeErrorf(pos, "pl.switchExpr",
			"only integer is support for swithc now, got %s", exprRef)
		return nil
	}
	// m := make(map[string]bool)

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
		pos := ast.ExprPos(c.Expr)
		b.CodeErrorf(pos, "pl.caseExpr",
			"nil cannot be used as expression for case")
		return nil
	}
	exprRef := e.R()
	if !exprRef.IsSingle() {
		pos := ast.ExprPos(c.Expr)
		b.CodeErrorf(pos, "pl.caseExpr",
			"expect single expression for case, got %s", exprRef)
		return nil
	}
	if !(types.IsInteger(exprRef.Type()) || types.IsConst(exprRef.Type())) {
		pos := ast.ExprPos(c.Expr)
		b.CodeErrorf(pos, "pl.caseExpr",
			"only integer is support for case now, got %s", exprRef)
		return nil
	}
	var stmts []tast.Stmt
	for _, stmt := range c.Stmts {
		s := b.buildStmt(stmt)
		if s != nil {
			stmts = append(stmts, s)
		}
	}
	return &tast.Case{Expr: e, Stmts: stmts, Fallthrough: c.Fallthrough == nil}
}
