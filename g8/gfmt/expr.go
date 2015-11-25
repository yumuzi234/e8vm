package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printExprs(p *fmt8.Printer, args ...interface{}) {
	for _, arg := range args {
		printExpr(p, arg)
	}
}

func printExpr(p *fmt8.Printer, expr ast.Expr) {
	switch expr := expr.(type) {
	case string:
		fmt.Fprintf(p, expr)
	case *ast.Operand:
		fmt.Fprintf(p, expr.Token.Lit)
	case *ast.OpExpr:
		if expr.A == nil {
			printExprs(p, expr.Op.Lit, expr.B)
		} else {
			printExprs(p, expr.A, " ", expr.Op.Lit, " ", expr.B)
		}
	case *ast.StarExpr:
		printExprs(p, "*", expr.Expr)
	case *ast.ParenExpr:
		printExprs(p, "(", expr.Expr, ")")
	case *ast.ExprList:
		for i, e := range expr.Exprs {
			if i != 0 {
				printExprs(p, ", ")
			}
			printExprs(p, e)
		}
	case *ast.CallExpr:
		if expr.Args != nil {
			printExprs(p, expr.Func, "(", expr.Args, ")")
		} else {
			printExprs(p, expr.Func, "()")
		}
	case *ast.IndexExpr:
		printExprs(p, expr.Array, "[", expr.Index, "]")
	case *ast.ArrayTypeExpr:
		if expr.Len == nil {
			printExprs(p, "[]", expr.Type)
		} else {
			printExprs(p, "[", expr.Len, "]", expr.Type)
		}
	case *ast.FuncTypeExpr:
		printExprs(p, "func")
	case *ast.MemberExpr:
		printExprs(p, expr.Expr, ".", expr.Sub.Lit)
	default:
		panic(fmt.Errorf("invalid expression type: %T", expr))
	}
}
