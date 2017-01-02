package gfmt

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func printExpr(f *formatter, expr interface{}) {
	switch expr := expr.(type) {
	case string:
		f.printStr(expr)
	case *lexing.Token:
		f.printToken(expr)
	case *ast.Operand:
		f.printToken(expr.Token)
	case *ast.OpExpr:
		if expr.A == nil {
			f.printExprs(expr.Op, expr.B)
		} else {
			f.printExprs(expr.A, " ", expr.Op, " ", expr.B)
		}
	case *ast.StarExpr:
		f.printExprs(expr.Star, expr.Expr)
	case *ast.ParenExpr:
		f.printExprs(expr.Lparen, expr.Expr, expr.Rparen)
	case *ast.ExprList:
		for i, e := range expr.Exprs {
			printExpr(f, e)
			if i < len(expr.Commas) {
				f.printExprs(expr.Commas[i], " ")
			}
		}
	case *ast.CallExpr:
		f.printExprs(expr.Func)
		f.printToken(expr.Lparen)
		if expr.Args != nil {
			printExprList(f, expr.Lparen, expr.Rparen, expr.Args)
		}
		f.printToken(expr.Rparen)
	case *ast.IndexExpr:
		if expr.Colon != nil {
			f.printExprs(expr.Array, expr.Lbrack)
			if expr.Index != nil {
				f.printExprs(expr.Index)
			}
			f.printExprs(expr.Colon)
			if expr.IndexEnd != nil {
				f.printExprs(expr.IndexEnd)
			}
			f.printExprs(expr.Rbrack)
		} else {
			f.printExprs(expr.Array, expr.Lbrack, expr.Index, expr.Rbrack)
		}
	case *ast.ArrayTypeExpr:
		if expr.Len == nil {
			f.printExprs(expr.Lbrack, expr.Rbrack, expr.Type)
		} else {
			f.printExprs(expr.Lbrack, expr.Len, expr.Rbrack, expr.Type)
		}
	case *ast.FuncTypeExpr:
		f.printToken(expr.Kw)
		printFuncSig(f, expr.FuncSig)
	case *ast.MemberExpr:
		f.printExprs(expr.Expr, expr.Dot, expr.Sub)
	case *ast.ArrayLiteral:
		if expr.Type.Len != nil {
			f.printExprs(expr.Type.Lbrack, expr.Type.Len, expr.Type.Rbrack)
		} else {
			f.printExprs(expr.Type.Lbrack, expr.Type.Rbrack)
		}

		f.printExprs(expr.Type.Type)
		if expr.Exprs != nil {
			f.printToken(expr.Lbrace)
			printExprList(f, expr.Lbrace, expr.Rbrace, expr.Exprs)
			f.printToken(expr.Rbrace)
		} else {
			f.printExprs(expr.Lbrace, expr.Rbrace)
		}
	default:
		f.errorf(nil, "invalid expression type: %T", expr)
	}
}
