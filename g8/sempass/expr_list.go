package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildExprList(b *Builder, list *ast.ExprList) *tast.ExprList {
	ret := tast.NewExprList()
	if list == nil {
		return ret
	}
	n := list.Len()
	if n == 0 {
		return ret
	}
	for _, expr := range list.Exprs {
		ex := b.BuildExpr(expr)
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
