package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseIndexExpr(p *parser, lead ast.Expr) *ast.IndexExpr {
	if !p.SeeOp("[") {
		panic("parseIndexExpr() must start with '['")
	}

	ret := new(ast.IndexExpr)

	ret.Array = lead
	ret.Lbrack = p.Shift()

	if !p.SeeOp(":") {
		ret.Index = p.parseExpr()
		if ret.Index == nil {
			return nil
		}
	}

	if p.SeeOp(":") {
		ret.Colon = p.Shift()
		if !p.SeeOp("]") {
			ret.IndexEnd = p.parseExpr()
			if ret.IndexEnd == nil {
				return nil
			}
		}
	}

	ret.Rbrack = p.ExpectOp("]")

	return ret
}
