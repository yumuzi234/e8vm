package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildExpr(b *Builder, expr ast.Expr) tast.Expr {
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
	case *ast.CallExpr:
		panic("todo")
	case *ast.IndexExpr:
		panic("todo")
	case *ast.ArrayTypeExpr:
		t := b.BuildType(expr)
		return &tast.Type{tast.NewRef(&types.Type{t})}
	case *ast.FuncTypeExpr:
		t := b.BuildType(expr)
		return &tast.Type{tast.NewRef(&types.Type{t})}
	}

	b.Errorf(ast.ExprPos(expr), "invalid or not implemented: %T", expr)
	return nil
}

func buildConstExpr(b *Builder, expr ast.Expr) tast.Expr {
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
