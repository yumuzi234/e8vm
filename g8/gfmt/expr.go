package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

func printExprs(f *formatter, args ...interface{}) {
	for _, arg := range args {
		printExpr(f, arg)
	}
}

func printExpr(f *formatter, expr ast.Expr) {
	switch expr := expr.(type) {
	case string:
		f.printStr(expr)
	case *lex8.Token:
		f.printToken(expr)
	case *ast.Operand:
		f.printToken(expr.Token)
	case *ast.OpExpr:
		if expr.A == nil {
			printExprs(f, expr.Op, expr.B)
		} else {
			printExprs(f, expr.A, " ", expr.Op, " ", expr.B)
		}
	case *ast.StarExpr:
		printExprs(f, expr.Star, expr.Expr)
	case *ast.ParenExpr:
		printExprs(f, expr.Lparen, expr.Expr, expr.Rparen)
	case *ast.ExprList:
		for i, e := range expr.Exprs {
			if i != 0 {
				printExprs(f, expr.Commas[i-1], " ")
			}
			printExpr(f, e)
		}
	case *ast.CallExpr:
		if expr.Args != nil {
			printExprs(f, expr.Func, expr.Lparen, expr.Args, expr.Rparen)
		} else {
			printExprs(f, expr.Func, expr.Lparen, expr.Rparen)
		}
	case *ast.IndexExpr:
		printExprs(f, expr.Array, expr.Lbrack, expr.Index, expr.Rbrack)
	case *ast.ArrayTypeExpr:
		if expr.Len == nil {
			printExprs(f, expr.Lbrack, expr.Rbrack, expr.Type)
		} else {
			printExprs(f, expr.Lbrack, expr.Len, expr.Rbrack, expr.Type)
		}
	case *ast.FuncTypeExpr:
		f.printToken(expr.Kw)
	case *ast.MemberExpr:
		printExprs(f, expr.Expr, expr.Dot, expr.Sub)
	default:
		panic(fmt.Errorf("invalid expression type: %T", expr))
	}
}
