package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

func buildConstExpr(b *builder, expr ast.Expr) *ref {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		return buildConstOperand(b, expr)
	case *ast.MemberExpr:
		return buildConstMember(b, expr)
	case *ast.OpExpr:
		return buildConstOpExpr(b, expr)
	case *ast.ParenExpr:
		return buildConstExpr(b, expr.Expr)
	default:
		b.Errorf(ast.ExprPos(expr), "expect a const expression")
		return nil
	}
}

func buildExpr(b *builder, expr ast.Expr) *ref {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		return buildOperand(b, expr)
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
		t := buildArrayType(b, expr)
		if t == nil {
			return nil
		}
		return newTypeRef(t)
	case *ast.FuncTypeExpr:
		t := buildFuncType(b, nil, expr.FuncSig)
		if t == nil {
			return nil
		}
		return newTypeRef(t)
	default:
		b.Errorf(ast.ExprPos(expr), "invalid or not implemented: %T", expr)
		return nil
	}
}

func buildExprStmt(b *builder, expr ast.Expr) {
	if e, ok := expr.(*ast.CallExpr); ok {
		buildCallExpr(b, e)
	} else {
		b.Errorf(ast.ExprPos(expr), "invalid expression statement")
	}
}
