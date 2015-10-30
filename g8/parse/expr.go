package parse

import (
	"io"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

func parseExpr(p *parser) ast.Expr {
	return parseBinaryExpr(p, 0)
}

// Exprs parses a list of expressions and returns an array of ast node of
// these expressions.
func Exprs(f string, r io.Reader) ([]ast.Expr, []*lex8.Error) {
	var ret []ast.Expr

	p, _ := newParser(f, r, false)
	p.exprFunc = parseExpr

	for !p.See(lex8.EOF) {
		expr := p.parseExpr()
		if expr != nil {
			ret = append(ret, expr)
		}

		p.ExpectSemi()
		if p.InError() {
			p.skipErrStmt()
		}
	}

	if es := p.Errs(); es != nil {
		return nil, es
	}

	return ret, nil
}
