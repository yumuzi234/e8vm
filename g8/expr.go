package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

// to replace buildExpr in the future
func buildExpr2(b *builder, expr tast.Expr) *ref {
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
		from := buildExpr2(b, expr.From)
		return buildCast(b, from, expr.T)
	case *tast.MemberExpr:
		return buildMember(b, expr)
	case *tast.OpExpr:
		return buildOpExpr(b, expr)
	case *tast.StarExpr:
		return buildStarExpr(b, expr)
	case *tast.CallExpr:
		return buildCallExpr(b, expr)
	case *tast.IndexExpr:
		return genIndexExpr(b, expr)
	case *tast.ExprList:
		return buildExprList(b, expr)
	}
	panic(fmt.Errorf("buildExpr2 not implemented for %T", expr))
}

func buildConstExpr(b *builder, expr ast.Expr) *ref {
	c := b.spass.BuildConstExpr(expr)
	if c == nil {
		return nil
	}
	return buildConst(b, c)
}

func buildExpr(b *builder, expr ast.Expr) *ref {
	e := b.spass.BuildExpr(expr)
	if e == nil {
		return nil
	}
	return buildExpr2(b, e)
}

func buildExprStmt(b *builder, expr ast.Expr) {
	if e, ok := expr.(*ast.CallExpr); ok {
		buildExpr(b, e)
	} else {
		b.Errorf(ast.ExprPos(expr), "invalid expression statement")
	}
}
