package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func seeType(p *parser) bool {
	if p.See(Ident) {
		return true
	}
	if p.SeeOp("*", "[", "(") {
		return true
	}
	if p.SeeKeyword("func") {
		return true
	}
	return false
}

func parseType(p *parser) ast.Expr {
	if p.See(Ident) {
		ret := ast.NewOperand(p.Shift())
		if p.SeeOp(".") {
			return parseMemberExpr(p, ret)
		}
		return ret
	} else if p.SeeOp("*") {
		ret := new(ast.StarExpr)
		ret.Star = p.Shift()
		ret.Expr = p.parseType()
		return ret
	} else if p.SeeOp("[") {
		ret := new(ast.ArrayTypeExpr)
		ret.Lbrack = p.Shift()
		if !p.SeeOp("]") {
			ret.Len = p.parseExpr()
			if ret.Len == nil {
				return nil
			}
		}
		ret.Rbrack = p.ExpectOp("]")
		ret.Type = p.parseType()
		return ret
	} else if p.SeeOp("(") {
		ret := new(ast.ParenExpr)
		ret.Lparen = p.Shift()
		ret.Expr = p.parseType()
		ret.Rparen = p.ExpectOp(")")
		return ret
	} else if p.SeeKeyword("func") {
		ret := new(ast.FuncTypeExpr)
		ret.Kw = p.Shift()
		ret.FuncSig = parseFuncSig(p)
		return ret
	}

	tok := p.Token()
	p.CodeErrorfHere("pl.expectType",
		"expect a type, got %s %q", p.TypeStr(tok.Type), tok.Lit)
	return nil
}
