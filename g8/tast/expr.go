package tast

import (
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// This is the this pointer.
type This struct{ *Ref }

// Const is a constant.
type Const struct{ *Ref }

// Type is a type expression
type Type struct{ *Ref }

// NewType creates a new type expression of a particular type.
func NewType(t types.T) *Type {
	return &Type{NewTypeRef(t)}
}

// Cast cast from one type of reference to another
type Cast struct {
	From Expr
	*Ref
}

// Ident is an identifier.
type Ident struct {
	Token *lex8.Token
	*Ref
	Symbol *sym8.Symbol
}

// MemberExpr is an expression of "a.b"
type MemberExpr struct {
	Expr Expr
	Sub  *lex8.Token
	*Ref
	Symbol *sym8.Symbol
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
	Array, Index, IndexEnd Expr
	HasColon               bool
	*Ref
}

// ExprList is a list of expressions.
type ExprList struct {
	Exprs []Expr
	*Ref
}

// Len returns the length of the expression list.
func (lst *ExprList) Len() int {
	return len(lst.Exprs)
}

// NewExprList creates a new expression list.
func NewExprList() *ExprList {
	return &ExprList{Ref: Void}
}

// Append appends an expression into the expression list.
func (lst *ExprList) Append(expr Expr) {
	ref := expr.R()
	lst.Ref = AppendRef(lst.Ref, ref)
	lst.Exprs = append(lst.Exprs, expr)
}
