package tast

import (
	"fmt"

	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

// Operand is an operand expression
type Operand struct {
	*lex8.Token
	T types.T
}

// MemberExpr is an expression of "a.b"
type MemberExpr struct {
	Expr Expr
	Sub  *lex8.Token
	T    types.T
}

// OpExpr is an expression likfe "a+b"
type OpExpr struct {
	A  Expr
	Op *lex8.Token
	B  Expr
	T  types.T
}

// StarExpr is an expression like "*a"
type StarExpr struct {
	Expr Expr
	T    types.T
}

// CallExpr is an expression like "f(x)"
type CallExpr struct {
	Func Expr
	Args *ExprList
	T    types.T
}

// IndexExpr is an expression like "a[b:c]"
// Both b and c are optional.
type IndexExpr struct {
	Array    Expr
	Index    Expr
	IndexEnd Expr
	T        types.T
}

// ArrayTypeExpr is an expresion like "[x]b".
// x is optional.
type ArrayTypeExpr struct {
	Len  Expr
	Type Expr
	T    types.T
}

// Para is a function parameter.
type Para struct {
	Ident *lex8.Token
	T     types.T
}

// FuncTypeExpr is an expression like "func f(t a)".
type FuncTypeExpr struct {
	Args []*Para
	Rets []*Para
	T    types.T
}

// ExprList is a list of expressions
type ExprList struct {
	Exprs []Expr
}

// ExprType returns the type of an expression node.
func ExprType(expr Expr) types.T {
	switch expr := expr.(type) {
	case *Operand:
		return expr.T
	case *MemberExpr:
		return expr.T
	case *OpExpr:
		return expr.T
	case *StarExpr:
		return expr.T
	case *CallExpr:
		return expr.T
	case *IndexExpr:
		return expr.T
	case *ArrayTypeExpr:
		return expr.T
	case *Para:
		return expr.T
	case *FuncTypeExpr:
		return expr.T
	default:
		panic(fmt.Errorf("invalid tast expr node: %T", expr))
	}
}
