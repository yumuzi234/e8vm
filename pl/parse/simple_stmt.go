package parse

import (
	"e8vm.io/e8vm/pl/ast"
)

func parseSimpleStmtOrExpr(p *parser, needSemi bool) (ast.Stmt, ast.Expr) {
	if p.See(Semi) {
		ret := new(ast.EmptyStmt)
		ret.Semi = p.Shift()
		return ret, nil
	}

	exprs := parseExprList(p)
	if exprs == nil {
		p.Next() // always make some progress
		return nil, nil
	} else if p.SeeOp(
		"=", "+=", "-=", "*=", "/=", "%=",
		"&=", "|=", "^=", "<<=", ">>=",
	) {
		// assigns statement
		ret := new(ast.AssignStmt)
		ret.Left = exprs
		ret.Assign = p.Shift()
		ret.Right = parseExprList(p)
		if needSemi {
			ret.Semi = p.ExpectSemi()
		}
		return ret, nil
	} else if p.SeeOp(":=") {
		// define statement
		ret := new(ast.DefineStmt)
		ret.Left = exprs
		ret.Define = p.Shift()
		ret.Right = parseExprList(p)
		if needSemi {
			ret.Semi = p.ExpectSemi()
		}
		return ret, nil
	} else if p.SeeOp("++", "--") {
		ret := new(ast.IncStmt)
		op := p.Shift()
		if exprs.Len() != 1 {
			p.ErrorfHere("%s on expression list", op.Lit)
		} else {
			ret.Expr = exprs.Exprs[0]
		}
		ret.Op = op
		if needSemi {
			ret.Semi = p.ExpectSemi()
		}
		return ret, nil
	}

	if exprs.Len() != 1 {
		p.ErrorfHere("expect expression, but got a list")
		p.BailOut()
		return nil, nil
	}

	expr := exprs.Exprs[0]

	if !needSemi {
		ret := new(ast.ExprStmt)
		ret.Expr = expr
		return ret, nil
	} else if semi := p.AcceptSemi(); semi != nil {
		ret := new(ast.ExprStmt)
		ret.Expr = expr
		ret.Semi = semi
		return ret, nil
	}

	return nil, expr
}

func parseSimpleStmt(p *parser) ast.Stmt {
	ret, expr := parseSimpleStmtOrExpr(p, true)
	if expr != nil {
		// semi is missing
		p.ExpectSemi()
		return nil
	}
	return ret
}

func parseSimpleStmtNoSemi(p *parser) ast.Stmt {
	ret, expr := parseSimpleStmtOrExpr(p, false)
	if expr != nil {
		panic("bug")
	}
	return ret
}
