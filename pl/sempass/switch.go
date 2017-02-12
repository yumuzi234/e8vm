package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildSwitchStmt(b *builder, stmt *ast.SwitchStmt) tast.Stmt {
	e := buildSwitchExpr(b, stmt.Expr)
	var cases []*tast.Case
	m := make(map[int64][]ast.Expr)
	for _, c := range stmt.Cases {
		ret := buildCase(b, m, c, e.Type())
		cases = append(cases, ret)
	}

	for v, exprs := range m {
		if len(exprs) > 1 {
			for _, e := range exprs {
				b.CodeErrorf(ast.ExprPos(e), "pl.caseExpr.dulplicated",
					"dulplicated case const, %d", v)
			}
		}
	}
	return &tast.SwitchStmt{Expr: e, Cases: cases}
}

func buildSwitchExpr(b *builder, expr ast.Expr) tast.Expr {
	e := b.buildExpr(expr)
	if e == nil {
		return nil
	}
	exprRef := e.R()
	if !exprRef.IsSingle() {
		pos := ast.ExprPos(expr)
		b.CodeErrorf(pos, "pl.switchExpr.notSingle",
			"expect single expression for switch, got %s", exprRef)
		return nil
	}
	if !(types.IsInteger(exprRef.Type()) || types.IsConst(exprRef.Type())) {
		pos := ast.ExprPos(expr)
		b.CodeErrorf(pos, "pl.swithExpr.notSupport",
			"only integer is support for case now, got %s", exprRef)
		return nil
	}
	return e
}

func buildCase(b *builder, m map[int64][]ast.Expr,
	c *ast.Case, t types.T) *tast.Case {
	var e tast.Expr
	if c.Kw.Lit == "case" {
		e = buildCaseExpr(b, m, c, t)
	}
	var stmts []tast.Stmt
	b.scope.Push()
	defer scopePopAndCheck(b)

	for _, stmt := range c.Stmts {
		s := b.buildStmt(stmt)
		if s != nil {
			stmts = append(stmts, s)
		}
	}
	return &tast.Case{Expr: e, Stmts: stmts, Fallthrough: c.Fallthrough == nil}
}

func buildCaseExpr(b *builder, m map[int64][]ast.Expr,
	c *ast.Case, t types.T) tast.Expr {
	e := b.buildExpr(c.Expr)
	pos := ast.ExprPos(c.Expr)
	if e == nil {
		b.CodeErrorf(pos, "pl.caseExpr.notConst",
			"only const integer value is allowed for case, got nil")
		return nil
	}
	r := e.R()
	if !r.IsSingle() {
		b.CodeErrorf(pos, "pl.caseExpr.notSingle",
			"expect single expression for case, got %s", r)
		return nil
	}
	v, ok := types.NumConst(r.Type())
	if !ok {
		b.CodeErrorf(pos, "pl.caseExpr.notConst",
			"only const integer value is allowed for case, got %s", r)
		return nil
	}
	if !types.IsConst(t) {
		e = constCast(b, pos, v, e, t)
	}
	m[v] = append(m[v], c.Expr)

	return e
}
