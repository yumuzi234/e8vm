package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func appendRef(base, toAdd *tast.Ref) *tast.Ref {
	if base == nil {
		return toAdd
	}
	if base.List == nil {
		return &tast.Ref{
			List: []*tast.Ref{base, toAdd},
		}
	}

	base.List = append(base.List, toAdd)
	return base
}

func buildExprList(b *builder, list *ast.ExprList) *tast.ExprList {
	ret := new(tast.ExprList)
	if list == nil {
		return ret
	}
	n := list.Len()
	if n == 0 {
		return ret
	}

	for _, expr := range list.Exprs {
		ex := b.buildExpr(expr)
		if ex == nil {
			return nil
		}

		ref := tast.ExprRef(ex)
		if ref.List != nil {
			b.Errorf(ast.ExprPos(expr), "cannot composite list in a list")
			return nil
		}
		ret.Ref = appendRef(ret.Ref, ref)
		ret.Exprs = append(ret.Exprs, ex)
	}
	return ret
}
