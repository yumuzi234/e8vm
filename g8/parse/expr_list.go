package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

func parseExprListClosed(p *parser, closeWith string) *ast.ExprList {
	ret := new(ast.ExprList)

	for {
		expr := p.parseExpr()
		if expr == nil {
			return nil
		}
		ret.Exprs = append(ret.Exprs, expr)

		if p.SeeOp(closeWith) {
			return ret
		}

		comma := p.ExpectOp(",")
		if comma == nil {
			return nil
		}
		ret.Commas = append(ret.Commas, comma)

		// could be a trailing comma
		if p.SeeOp(closeWith) {
			return ret
		}
	}

	return ret
}

func parseExprList(p *parser) *ast.ExprList {
	ret := new(ast.ExprList)
	for {
		expr := p.parseExpr()
		if expr == nil {
			return nil
		}
		ret.Exprs = append(ret.Exprs, expr)
		if !p.SeeOp(",") {
			break
		}
		ret.Commas = append(ret.Commas, p.Shift())
	}
	return ret
}
