package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

func parseMemberExpr(p *parser, lead ast.Expr) *ast.MemberExpr {
	if !p.SeeOp(".") {
		panic("parseMemberExpr() must start with '.'")
	}

	return &ast.MemberExpr{
		Expr: lead,
		Dot:  p.Shift(),
		Sub:  p.Expect(Ident),
	}
}
