package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

// for <cond> { <stmts> }
func parseForStmt(p *parser) *ast.ForStmt {
	if !p.SeeKeyword("for") {
		panic("must start with keyword")
	}

	ret := new(ast.ForStmt)
	ret.Kw = p.Shift()
	if !p.SeeOp("{") {
		stmt, expr := parseSimpleStmtOrExpr(p, true)
		if stmt != nil {
			ret.ThreeFold = true

			ret.Init = stmt
			if !p.SeeSemi() {
				ret.Cond = parseExpr(p)
			}

			p.ExpectSemi()
			if !p.InError() {
				ret.Iter = parseSimpleStmtNoSemi(p)
			}
		} else if expr != nil {
			ret.Cond = expr
		}

		if p.InError() {
			return ret
		}
	}

	ret.Body = parseBlock(p)
	if p.InError() {
		return ret
	}

	ret.Semi = p.ExpectSemi()
	return ret
}
