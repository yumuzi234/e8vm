package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
)

func buildIdentExprList(b *builder, list *ast.ExprList) (
	idents []*lex8.Token, firstError ast.Expr,
) {
	ret := make([]*lex8.Token, 0, list.Len())
	for _, expr := range list.Exprs {
		op, ok := expr.(*ast.Operand)
		if !ok {
			return nil, expr
		}
		if op.Token.Type != parse.Ident {
			return nil, expr
		}

		ret = append(ret, op.Token)
	}

	return ret, nil
}

func buildExprList(b *builder, list *tast.ExprList) *ref {
	if list == nil {
		return new(ref)
	}
	n := list.Len()
	if n == 0 {
		return new(ref)
	} else if n == 1 {
		return b.buildExpr(list.Exprs[0])
	}

	ret := new(ref)
	for _, expr := range list.Exprs {
		ret = appendRef(ret, b.buildExpr(expr))
	}
	return ret
}
