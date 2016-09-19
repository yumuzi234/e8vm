package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseReturnStmt(p *parser, withSemi bool) *ast.ReturnStmt {
	ret := new(ast.ReturnStmt)
	ret.Kw = p.ExpectKeyword("return")
	if !p.SeeSemi() {
		ret.Exprs = parseExprList(p)
	}
	if withSemi {
		ret.Semi = p.ExpectSemi()
	}
	return ret
}
