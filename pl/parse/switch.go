package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseSwitchStmt(p *parser) *ast.SwitchStmt {
	if !p.SeeKeyword("switch") {
		panic("must start with keyword switch")
	}
	ret := new(ast.SwitchStmt)
	ret.Kw = p.Shift()
	ret.Expr = p.parseExpr()
	if p.InError() {
		return ret
	}
	if !p.SeeOp("{") {
		p.CodeErrorfHere("missingSwitchBody",
			"missing switch body, need '{'")
		return ret
	}
	ret.Lbrace = p.Shift()
	ret.Cases = parseCases(p)
	ret.Rbrace = p.ExpectOp("}")
	ret.Semi = p.ExpectSemi()
	return ret
}
