package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

func printExprs(p *fmt8.Printer, m *matcher, args ...interface{}) {
	for _, arg := range args {
		printExpr(p, m, arg)
	}
}

func printExpr(p *fmt8.Printer, m *matcher, expr ast.Expr) {
	switch expr := expr.(type) {
	case string:
		fmt.Fprintf(p, expr)
	case *lex8.Token:
		printToken(p, m, expr)
	case *ast.Operand:
		printToken(p, m, expr.Token)
	case *ast.OpExpr:
		if expr.A == nil {
			printExprs(p, m, expr.Op, expr.B)
		} else {
			printExprs(p, m, expr.A, " ", expr.Op, " ", expr.B)
		}
	case *ast.StarExpr:
		printExprs(p, m, expr.Star, expr.Expr)
	case *ast.ParenExpr:
		printExprs(p, m, expr.Lparen, expr.Expr, expr.Rparen)
	case *ast.ExprList:
		for i, e := range expr.Exprs {
			if i != 0 {
				printExprs(p, m, expr.Commas[i-1], " ")
			}
			printExpr(p, m, e)
		}
	case *ast.CallExpr:
		if expr.Args != nil {
			printExprs(p, m, expr.Func, expr.Lparen, expr.Args, expr.Rparen)
		} else {
			printExprs(p, m, expr.Func, expr.Lparen, expr.Rparen)
		}
	case *ast.IndexExpr:
		printExprs(p, m, expr.Array, expr.Lbrack, expr.Index, expr.Rbrack)
	case *ast.ArrayTypeExpr:
		if expr.Len == nil {
			printExprs(p, m, expr.Lbrack, expr.Rbrack, expr.Type)
		} else {
			printExprs(p, m, expr.Lbrack, expr.Len, expr.Rbrack, expr.Type)
		}
	case *ast.FuncTypeExpr:
		printToken(p, m, expr.Kw)
	case *ast.MemberExpr:
		printExprs(p, m, expr.Expr, expr.Dot, expr.Sub)
	default:
		panic(fmt.Errorf("invalid expression type: %T", expr))
	}
}
