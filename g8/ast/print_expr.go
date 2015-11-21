package ast

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
)

func printExprs(p *fmt8.Printer, args ...interface{}) {
	for _, arg := range args {
		printExpr(p, arg)
	}
}

func printExpr(p *fmt8.Printer, expr Expr) {
	switch expr := expr.(type) {
	case string:
		fmt.Fprintf(p, expr)
	case *Operand:
		fmt.Fprintf(p, expr.Token.Lit)
	case *OpExpr:
		if expr.A == nil {
			printExprs(p, expr.Op.Lit, expr.B)
		} else {
			printExprs(p, expr.A, " ", expr.Op.Lit, " ", expr.B)
		}
	case *StarExpr:
		printExprs(p, "*", expr.Expr)
	case *ParenExpr:
		printExprs(p, "(", expr.Expr, ")")
	case *ExprList:
		for i, e := range expr.Exprs {
			if i != 0 {
				printExprs(p, ", ")
			}
			printExprs(p, e)
		}
	case *CallExpr:
		if expr.Args != nil {
			printExprs(p, expr.Func, "(", expr.Args, ")")
		} else {
			printExprs(p, expr.Func, "()")
		}
	case *IndexExpr:
		printExprs(p, expr.Array, "[", expr.Index, "]")
	case *ArrayTypeExpr:
		if expr.Len == nil {
			printExprs(p, "[]", expr.Type)
		} else {
			printExprs(p, "[", expr.Len, "]", expr.Type)
		}
	case *FuncTypeExpr:
		printExprs(p, "func")
	case *MemberExpr:
		printExprs(p, expr.Expr, ".", expr.Sub.Lit)
	default:
		panic(fmt.Errorf("invalid expression type: %T", expr))
	}
}
