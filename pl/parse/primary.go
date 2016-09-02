package parse

import (
	"e8vm.io/e8vm/pl/ast"
)

func parsePrimaryExpr(p *parser) ast.Expr {
	ret := parseOperand(p)
	if ret == nil {
		return nil
	}

	for {
		if p.SeeOp("(") {
			ret = parseCallExpr(p, ret)
		} else if p.SeeOp("[") {
			ret = parseIndexExpr(p, ret)
		} else if p.SeeOp(".") {
			ret = parseMemberExpr(p, ret)
		} else {
			break
		}

		if ret == nil {
			return nil
		}
	}

	return ret
}
