package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

func parseCallExpr(p *parser, lead ast.Expr) *ast.CallExpr {
	if !p.SeeOp("(") {
		panic("parseCallExpr() must start with '('")
	}

	lp := p.Shift()

	if p.SeeOp(")") {
		// no args
		rp := p.Shift()
		return &ast.CallExpr{
			Func:   lead,
			Args:   nil,
			Lparen: lp,
			Rparen: rp,
		}
	}

	lst := parseExprListClosed(p, ")")
	if p.InError() {
		return nil
	}
	rp := p.ExpectOp(")")
	if rp == nil {
		return nil
	}

	return &ast.CallExpr{
		Func:   lead,
		Args:   lst,
		Lparen: lp,
		Rparen: rp,
	}
}
