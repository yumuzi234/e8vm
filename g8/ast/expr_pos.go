package ast

import (
	"fmt"

	"e8vm.io/e8vm/lex8"
)

// ExprPos returns the starting position of an expression.
func ExprPos(e Expr) *lex8.Pos {
	switch e := e.(type) {
	case *Operand:
		return e.Token.Pos
	case *OpExpr:
		if e.A != nil {
			return ExprPos(e.A)
		}
		return e.Op.Pos
	case *ParenExpr:
		return e.Lparen.Pos
	case *ExprList:
		if len(e.Exprs) == 0 {
			return nil
		}
		return ExprPos(e.Exprs[0])
	case *CallExpr:
		return ExprPos(e.Func)
	case *StarExpr:
		return e.Star.Pos
	case *IndexExpr:
		return ExprPos(e.Array)
	case *MemberExpr:
		return ExprPos(e.Expr)
	case *ArrayTypeExpr:
		return e.Lbrack.Pos
	case *ArrayLiteral:
		return ExprPos(e.Type)
	case *FuncTypeExpr:
		return e.Kw.Pos
	default:
		panic(fmt.Errorf("invalid expression type: %T", e))
	}
}
