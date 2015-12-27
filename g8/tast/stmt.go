package tast

import (
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// ExprStmt is a statement with just an expression.
type ExprStmt struct {
	Expr
}

// AssignStmt is an assignment statement, like "a,b=x,y".
type AssignStmt struct {
	Left  *ExprList
	Right Expr
}

// DefineStmt is a define statement, like "a,b:=x,y".
type DefineStmt struct {
	Left  []*sym8.Symbol
	Right Expr // zero out Left if Right==nil
}

// Block is a statement block; it is also a statement itself.
type Block struct {
	Stmts []Stmt
}

// IncStmt is an "i++" or "i--".
type IncStmt struct {
	Expr Expr
	Op   *lex8.Token
}

// ContinueStmt is a "continue"
type ContinueStmt struct{}

// BreakStmt is a "break"
type BreakStmt struct{}

// ReturnStmt is a statement like "return a,b"
type ReturnStmt struct {
	Exprs Expr
}

// IfStmt is an if statement.
type IfStmt struct {
	Expr Expr
	Body Stmt
	Else *ElseStmt
}

// ElseStmt the else stmt that follows an if statement.
type ElseStmt struct {
	Expr Expr // optional for else if
	Body *Block
	Next *ElseStmt // next else statement
}

// ForStmt is a for loop statement.
type ForStmt struct {
	ThreeFold bool
	Init      Stmt
	Cond      Expr
	Iter      Stmt
	Body      *Block
}
