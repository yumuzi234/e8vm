package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildExpr(b *builder, expr ast.Expr) tast.Expr {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		return buildOperand(b, expr)
	case *ast.ParenExpr:
		return buildExpr(b, expr.Expr)
	case *ast.MemberExpr:
		return buildMember(b, expr)
	case *ast.OpExpr:
		return buildOpExpr(b, expr)
	case *ast.StarExpr:
		return buildStarExpr(b, expr)
	case *ast.IndexExpr:
		return buildIndexExpr(b, expr)
	case *ast.CallExpr:
		return buildCallExpr(b, expr)
	case *ast.ArrayTypeExpr:
		t := b.BuildType(expr)
		if t == nil {
			return nil
		}
		return tast.NewType(t)
	case *ast.FuncTypeExpr:
		t := b.BuildType(expr)
		if t == nil {
			return nil
		}
		return tast.NewType(t)
	case *ast.ExprList:
		return buildExprList(b, expr)
	}

	b.Errorf(ast.ExprPos(expr), "invalid or not implemented: %T", expr)
	return nil
}

func buildConstExpr(b *builder, expr ast.Expr) tast.Expr {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return buildConstExpr(b, expr.Expr)
	case *ast.Operand:
		return buildConstOperand(b, expr)
	case *ast.MemberExpr:
		return buildConstMember(b, expr)
	case *ast.OpExpr:
		return buildConstOpExpr(b, expr)
	}

	b.Errorf(ast.ExprPos(expr), "expect a const expression")
	return nil
}

func buildExprStmt(b *builder, expr ast.Expr) tast.Stmt {
	if e, ok := expr.(*ast.CallExpr); ok {
		ret := buildExpr(b, e)
		if ret == nil {
			return nil
		}
		return &tast.ExprStmt{ret}
	}

	b.Errorf(ast.ExprPos(expr), "invalid expression statement")
	return nil
}
