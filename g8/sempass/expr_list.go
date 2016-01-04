package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildExprList(b *builder, list *ast.ExprList) tast.Expr {
	ret := tast.NewExprList()
	if list == nil {
		return ret
	}
	n := list.Len()
	if n == 0 {
		return ret
	}
	if n == 1 {
		return b.buildExpr(list.Exprs[0])
	}

	for _, expr := range list.Exprs {
		ex := b.buildExpr(expr)
		if ex == nil {
			return nil
		}

		ref := ex.R()
		if !ref.IsSingle() {
			b.Errorf(ast.ExprPos(expr), "cannot put %s in a list", ref)
			return nil
		}

		ret.Append(ex)
	}
	return ret
}

func buildConstExprList(b *builder, list *ast.ExprList) tast.Expr {
	n := list.Len()
	if n == 0 {
		b.Errorf(ast.ExprPos(list), "const expression list of zero length")
		return nil
	}
	if n == 1 {
		return b.buildConst(list.Exprs[0])
	}

	ret := tast.NewExprList()
	for _, expr := range list.Exprs {
		ex := b.buildConst(expr)
		if ex == nil {
			return nil
		}
		ref := ex.R()
		if !ref.IsSingle() {
			b.Errorf(ast.ExprPos(expr), "cannot put %s in a list", ref)
			return nil
		}
		ret.Append(ex)
	}
	return ret
}
