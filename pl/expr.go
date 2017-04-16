package pl

import (
	"fmt"

	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildExpr(b *builder, expr tast.Expr) *ref {
	switch expr := expr.(type) {
	case *tast.Const:
		return buildConst(b, expr)
	case *tast.Ident:
		return buildIdent(b, expr)
	case *tast.This:
		return b.this
	case *tast.Type:
		t := expr.Ref.T.(*types.Type)
		return newRef(t, nil)
	case *tast.Cast:
		from := buildExpr(b, expr.From)
		return buildCasts(b, from, expr.R(), expr.NeedCast)
	case *tast.MemberExpr:
		return buildMember(b, expr)
	case *tast.OpExpr:
		return buildOpExpr(b, expr)
	case *tast.StarExpr:
		return buildStarExpr(b, expr)
	case *tast.CallExpr:
		return buildCallExpr(b, expr)
	case *tast.IndexExpr:
		return buildIndexExpr(b, expr)
	case *tast.ExprList:
		return buildExprList(b, expr)
	}
	panic(fmt.Errorf("buildExpr not implemented for %T", expr))
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
