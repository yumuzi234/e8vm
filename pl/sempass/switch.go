package sempass

import (
	"shanhu.io/smlvm/lexing"
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
	var cases []*tast.Case
	m := make(map[int64]*lexing.Pos)
	// TO DO
	// v, ok := types.NumConst(exprRef.Type())
	// if ok{}
	for _, c := range Cases {
		ret := buildCase(b, m, c)
		if ret == nil {
			return nil
		}
		cases = append(cases, ret)
	}

	return &tast.SwitchStmt{Expr: e, Cases: cases}
}

func buildCase(b *builder, m map[int64]*lexing.Pos, c *ast.Case) *tast.Case {
	var e tast.Expr

	// c.Expr will be nil for the default:
	if c.Expr != nil {
		e = b.buildExpr(c.Expr)
		pos := ast.ExprPos(c.Expr)
		if e == nil {
			return nil
		}
		r := e.R()
		v, ok := types.NumConst(r.Type())
		if !(r.IsSingle() && ok) {
			b.CodeErrorf(pos, "pl.caseExpr.notConst",
				"only const integer value is allowed for case, got %s", r)
			return nil
		}
		if m[v] != nil {
			b.CodeErrorf(m[v], "pl.caseExpr.dulplicated",
				"dulplicated case const", r)
			b.CodeErrorf(pos, "pl.caseExpr.dulplicated",
				"dulplicated case const", r)
			return nil
		}
		m[v] = pos
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
