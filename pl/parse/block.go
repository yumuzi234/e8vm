package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseBlock(p *parser) *ast.Block {
	ret := new(ast.Block)
	ret.Lbrace = p.ExpectOp("{")
	if ret.Lbrace == nil {
		return ret
	}

	for !(p.SeeOp("}") || p.See(lexing.EOF)) {
		if stmt := p.parseStmt(); stmt != nil {
			ret.Stmts = append(ret.Stmts, stmt)
		}
		p.skipErrStmt()
	}

	ret.Rbrace = p.ExpectOp("}")
	return ret
}
