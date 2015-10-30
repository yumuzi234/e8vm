package parse

import (
	"e8vm.io/e8vm/g8/ast"
)

func parseUnaryExpr(p *parser) ast.Expr {
	if p.SeeOp("+", "-", "!", "^", "&") {
		t := p.Shift()
		expr := parseUnaryExpr(p)
		return &ast.OpExpr{A: nil, Op: t, B: expr}
	} else if p.SeeOp("*") {
		star := p.Shift()
		expr := parseUnaryExpr(p)
		return &ast.StarExpr{star, expr}
	}

	return parsePrimaryExpr(p)
}

func opPrec(op string) int {
	switch op {
	case "||":
		return 0
	case "&&":
		return 1
	case "==", "!=", "<", "<=", ">=", ">":
		return 2
	case "+", "-", "|", "^":
		return 3
	case "*", "%", "/", "<<", ">>", "&":
		return 4
	}
	return -1
}

func parseBinaryExpr(p *parser, prec int) ast.Expr {
	ret := parseUnaryExpr(p)
	if p.InError() {
		return nil
	}

	if p.See(Operator) {
		startPrec := opPrec(p.Token().Lit)
		for i := startPrec; i >= prec; i-- {
			for p.See(Operator) {
				if opPrec(p.Token().Lit) != i {
					break
				}

				op := p.Shift()
				bop := new(ast.OpExpr)
				bop.A = ret
				bop.Op = op
				bop.B = parseBinaryExpr(p, i+1)
				ret = bop
			}
		}
	}

	return ret
}
