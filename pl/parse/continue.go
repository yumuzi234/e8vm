package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseContinueStmt(p *parser, withSemi bool) *ast.ContinueStmt {
	ret := new(ast.ContinueStmt)
	ret.Kw = p.ExpectKeyword("continue")
	if p.See(Ident) {
		ret.Label = p.Expect(Ident)
	}
	if withSemi {
		ret.Semi = p.ExpectSemi()
	}
	return ret
}
