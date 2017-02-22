package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
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
			b.CodeErrorf(ast.ExprPos(expr), "pl.multiRefInExprList",
				`cannot put %s in a expresssion list, 
				only single reference is allowed`, ref)
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
		ret := b.buildConst(list.Exprs[0])
		if ret == nil {
			return nil
		}
		return ret
	}

	ret := tast.NewExprList()
	for _, expr := range list.Exprs {
		ex := b.buildConst(expr)
		if ex == nil {
			return nil
		}
		ref := ex.R()
		if !ref.IsSingle() {
			b.CodeErrorf(ast.ExprPos(expr), "pl.multiRefInExprList",
				`cannot put %s in a expresssion list, 
				only single reference is allowed`, ref)
			return nil
		}
		ret.Append(ex)
	}
	return ret
}
