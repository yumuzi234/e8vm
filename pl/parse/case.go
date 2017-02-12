package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseCases(p *parser) []*ast.Case {
	var ret []*ast.Case
	for !(p.SeeOp("}") || p.See(lexing.EOF)) {
		if c := parseCase(p); c != nil {
			ret = append(ret, c)
		}
		p.skipErrStmt()
	}
	return ret
}

func parseCase(p *parser) *ast.Case {
	ret := new(ast.Case)
	if p.SeeKeyword("case") {
		ret.Kw = p.Shift()
		ret.Expr = parseExpr(p)
		if ret.Expr == nil {
			return nil
		}
	} else if p.SeeKeyword("default") {
		ret.Kw = p.Shift()
	} else {
		p.CodeErrorfHere("pl.missingCaseInSwitch",
			"must start with keyword case/default in switch")
		return nil
	}
	ret.Colon = p.ExpectOp(":")
	if ret.Colon == nil {
		return nil
	}
	for !(p.SeeKeyword("case") || p.SeeKeyword("default") ||
		p.SeeOp("}") || p.See(lexing.EOF)) {
		if p.SeeKeyword("fallthrough") {
			break
		}
		if stmt := p.parseStmt(); stmt != nil {
			ret.Stmts = append(ret.Stmts, stmt)
		}
		p.skipErrStmt()
	}
	if p.SeeKeyword("fallthrough") {
		ret.Fallthrough = &ast.FallthroughStmt{
			Kw:   p.Shift(),
			Semi: p.ExpectSemi(),
		}

		if p.InError() {
			return ret
		}

		if !(p.SeeKeyword("case") || p.SeeKeyword("default")) {
			p.CodeErrorfHere("pl.invalidFallthrough",
				"fallthrough must be followed by new switch case")
			return nil
		}

	}
	return ret
}
