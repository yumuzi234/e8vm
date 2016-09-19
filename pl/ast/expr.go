package ast

import (
	"shanhu.io/smlvm/lexing"
)

// Operand is an operand expression
type Operand struct {
	*lexing.Token
}

// MemberExpr is an expression of form A.B
type MemberExpr struct {
	Expr Expr
	Dot  *lexing.Token
	Sub  *lexing.Token
}

// OpExpr is a binary or unary operation that uses an operator
type OpExpr struct {
	A  Expr
	Op *lexing.Token
	B  Expr
}

// StarExpr is an expression after a '*'
type StarExpr struct {
	Star *lexing.Token
	Expr Expr
}

// ParenExpr is an expression in a pair of parenthesis
type ParenExpr struct {
	Lparen *lexing.Token
	Expr
	Rparen *lexing.Token
}

// ExprList is a list of expressions
type ExprList struct {
	Exprs  []Expr
	Commas []*lexing.Token
}

// Len returns the length of the expression list
func (list *ExprList) Len() int {
	if list == nil {
		return 0
	}
	return len(list.Exprs)
}

// CallExpr is a function call expression
type CallExpr struct {
	Func   Expr
	Lparen *lexing.Token
	Args   *ExprList
	Rparen *lexing.Token
}

// IndexExpr is fetching an element in an array or slice
type IndexExpr struct {
	Array    Expr
	Lbrack   *lexing.Token
	Index    Expr
	Colon    *lexing.Token
	IndexEnd Expr
	Rbrack   *lexing.Token
}

// ArrayTypeExpr is the type expression of an array or a slice
type ArrayTypeExpr struct {
	Lbrack *lexing.Token
	Len    Expr // optional
	Rbrack *lexing.Token
	Type   Expr
}

// ArrayLiteral is an array or slice literal
type ArrayLiteral struct {
	Type   *ArrayTypeExpr
	Lbrace *lexing.Token
	Exprs  *ExprList
	Rbrace *lexing.Token
}

// FuncTypeExpr is the type expression of a function pointer
type FuncTypeExpr struct {
	Kw      *lexing.Token
	FuncSig *FuncSig
}
