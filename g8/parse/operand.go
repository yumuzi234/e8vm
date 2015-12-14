package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

func parseOperand(p *parser) ast.Expr {
	if p.See(Ident) || p.See(Int) || p.See(Float) ||
		p.See(String) || p.See(Char) {
		return &ast.Operand{p.Shift(), nil}
	} else if p.SeeKeyword("this") {
		return &ast.Operand{p.Shift(), nil}
	} else if p.SeeOp("(") {
		lp := p.Shift()
		expr := p.parseExpr()
		rp := p.ExpectOp(")")
		if rp == nil {
			return nil
		}

		return &ast.ParenExpr{
			Lparen: lp,
			Rparen: rp,
			Expr:   expr,
		}
	} else if p.SeeOp("*") || p.SeeOp("[") || p.SeeKeyword("func") {
		return p.parseType()
	}

	t := p.Token()
	p.Errorf(t.Pos, "expect an operand, got %s", p.typeStr(t))

	return nil
}
