package tast

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/syms"
)

// ExprStmt is a statement with just an expression.
type ExprStmt struct {
	Expr
}

// AssignStmt is an assignment statement, like "a,b=x,y".
type AssignStmt struct {
	Left  Expr
	Op    *lexing.Token
	Right Expr
}

// Define is a define statement, like "a,b:=x,y".
type Define struct {
	Left  []*syms.Symbol
	Right Expr // zero out Left if Right==nil
}

// Block is a statement block; it is also a statement itself.
type Block struct {
	Stmts []Stmt
}

// IncStmt is an "i++" or "i--".
type IncStmt struct {
	Expr Expr
	Op   *lexing.Token
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
	Else Stmt
}

// ForStmt is a for loop statement.
type ForStmt struct {
	ThreeFold bool
	Init      Stmt
	Cond      Expr
	Iter      Stmt
	Body      Stmt
}
