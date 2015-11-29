package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

func buildExprList(b *builder, list *ast.ExprList) *ref {
	ret := new(ref)
	if list == nil {
		return ret // empty ref, for void
	}

	n := list.Len()
	if n == 0 {
		return ret // empty ref
	} else if n == 1 {
		return b.buildExpr(list.Exprs[0])
	}

	for _, expr := range list.Exprs {
		var ref *ref
		ref = b.buildExpr(expr)

		if ref == nil {
			return nil
		}
		if !ref.IsSingle() {
			b.Errorf(ast.ExprPos(expr), "cannot composite list in a list")
			return nil
		}

		ret = appendRef(ret, ref)
	}

	return ret
}

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
