package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
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
	}
	panic(fmt.Errorf("genExpr not implemented for %T", expr))
}

func buildConstExpr(b *builder, expr ast.Expr) *ref {
	if expr == nil {
		panic("bug")
	}
	e := b.spass.BuildConstExpr(expr)
	if e == nil {
		return nil
	}

	switch e := e.(type) {
	case *tast.Const:
		return buildConst(b, e)
	case *tast.Ident:
		// TODO: why return ident here?
		return buildConstIdent(b, e)
	}
	panic("bug")
}

func buildExpr(b *builder, expr ast.Expr) *ref {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.MemberExpr:
		return buildMember(b, expr)
	case *ast.ParenExpr:
		return buildExpr(b, expr.Expr)
	case *ast.OpExpr:
		return buildOpExpr(b, expr)
	case *ast.StarExpr:
		return buildStarExpr(b, expr)
	case *ast.CallExpr:
		return buildCallExpr(b, expr)
	case *ast.IndexExpr:
		return buildIndexExpr(b, expr)
	case *ast.ArrayTypeExpr:
		t := b.spass.BuildType(expr)
		if t == nil {
			return nil
		}
		return newTypeRef(t)
	case *ast.FuncTypeExpr:
		t := b.spass.BuildType(expr)
		if t == nil {
			return nil
		}
		return newTypeRef(t)
	}

	e := b.spass.BuildExpr(expr)
	if e == nil {
		return nil
	}
	return buildExpr2(b, e)
}

func buildExprStmt(b *builder, expr ast.Expr) {
	if e, ok := expr.(*ast.CallExpr); ok {
		buildCallExpr(b, e)
	} else {
		b.Errorf(ast.ExprPos(expr), "invalid expression statement")
	}
}
