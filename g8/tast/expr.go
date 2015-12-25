package tast

import (
	"fmt"

	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// This is the this pointer.
type This struct {
	*Ref
}

// Const is a constant.
type Const struct {
	*Ref
}

// Ident is an identifier.
type Ident struct {
	*lex8.Token
	*Ref
	Symbol *sym8.Symbol
}

// MemberExpr is an expression of "a.b"
type MemberExpr struct {
	Expr Expr
	Sub  *lex8.Token
	*Ref
}

// OpExpr is an expression likfe "a+b"
type OpExpr struct {
	A  Expr
	Op *lex8.Token
	B  Expr
	*Ref
}

// StarExpr is an expression like "*a"
type StarExpr struct {
	Expr Expr
	*Ref
}

// CallExpr is an expression like "f(x)"
type CallExpr struct {
	Func Expr
	Args *ExprList
	*Ref
}

// IndexExpr is an expression like "a[b:c]"
// Both b and c are optional.
type IndexExpr struct {
	Array    Expr
	Index    Expr
	IndexEnd Expr
	*Ref
}

/*
// ArrayTypeExpr is an expresion like "[x]b".
// x is optional.
type ArrayTypeExpr struct {
	Len  Expr // optional
	Type Expr
	*Ref
}

// Para is a function parameter.
type Para struct {
	Ident *lex8.Token
	*Ref
}

// FuncTypeExpr is an expression like "func f(t a)".
type FuncTypeExpr struct {
	Args []*Para
	Rets []*Para
	*Ref
}
*/

// ExprList is a list of expressions.
type ExprList struct {
	Exprs []Expr
	*Ref
}

// ExprRef returns the reference of the expression.
func ExprRef(expr Expr) *Ref {
	switch expr := expr.(type) {
	case *This:
		return expr.Ref
	case *Ident:
		return expr.Ref
	case *Const:
		return expr.Ref
	case *MemberExpr:
		return expr.Ref
	case *OpExpr:
		return expr.Ref
	case *StarExpr:
		return expr.Ref
	case *CallExpr:
		return expr.Ref
	case *IndexExpr:
		return expr.Ref
	case *ExprList:
		return expr.Ref
	default:
		panic(fmt.Errorf("invalid tast expr node: %T", expr))
	}
}
