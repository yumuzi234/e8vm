package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseFallthroughStmt(p *parser) *ast.FallthroughStmt {
	ret := new(ast.FallthroughStmt)
	ret.Kw = p.ExpectKeyword("fallthrough")
	return ret
}
