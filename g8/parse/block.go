package parse

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

func parseBlock(p *parser) *ast.Block {
	ret := new(ast.Block)
	ret.Lbrace = p.ExpectOp("{")
	if ret.Lbrace == nil {
		return ret
	}

	for !(p.SeeOp("}") || p.See(lex8.EOF)) {
		if stmt := p.parseStmt(); stmt != nil {
			ret.Stmts = append(ret.Stmts, stmt)
		}
		p.skipErrStmt()
	}

	ret.Rbrace = p.ExpectOp("}")
	return ret
}
