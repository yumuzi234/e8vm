package ast

import (
	"shanhu.io/smlvm/lexing"
)

// ExprStmt is a statement with just an expression
type ExprStmt struct {
	Expr
	Semi *lexing.Token
}

// AssignStmt is an assignment statement:
// exprList = exprList
type AssignStmt struct {
	Left   *ExprList
	Assign *lexing.Token
	Right  *ExprList
	Semi   *lexing.Token
}

// DefineStmt is a statement that defines one or a list of variables.
// exprList := exprList
type DefineStmt struct {
	Left   *ExprList
	Define *lexing.Token
	Right  *ExprList
	Semi   *lexing.Token
}

// Block is a statement block
type Block struct {
	Lbrace *lexing.Token
	Stmts  []Stmt
	Rbrace *lexing.Token
}

// BlockStmt is a block statement
type BlockStmt struct {
	*Block
	Semi *lexing.Token
}

// IfStmt is an if statement, possibly with an else of else if
// following
type IfStmt struct {
	If   *lexing.Token
	Expr Expr
	Body Stmt
	Else *ElseStmt // optional for else or else if
	Semi *lexing.Token
}

// ElseStmt is the dangling statement block after if
type ElseStmt struct {
	Else *lexing.Token
	If   *lexing.Token // optional
	Expr Expr          // optional for else if
	Body *Block
	Next *ElseStmt // next else statement
}

// SwitchStmt is the swithc statement block
type SwitchStmt struct {
	Kw     *lexing.Token
	Expr   Expr
	Lbrace *lexing.Token
	Cases  []*Case
	Rbrace *lexing.Token
	Semi   *lexing.Token
}

// Case is the inset statement block in switch
// default is included here, Kw will determine it is case or default
type Case struct {
	Kw          *lexing.Token
	Expr        Expr
	Colon       *lexing.Token
	Stmts       []Stmt
	Fallthrough *FallthroughStmt
}

// ForStmt is a loop statement
type ForStmt struct {
	Kw        *lexing.Token
	ThreeFold bool
	Init      Stmt
	Cond      Expr
	CondSemi  *lexing.Token
	Iter      Stmt
	Body      *Block
	Semi      *lexing.Token
}

// ReturnStmt is a statement of return.
// return <expr>
type ReturnStmt struct {
	Kw    *lexing.Token
	Exprs *ExprList
	Semi  *lexing.Token
}

// IncStmt is an "i++" or "i--".
type IncStmt struct {
	Expr Expr
	Op   *lexing.Token
	Semi *lexing.Token
}

// ContinueStmt is the continue statement
// continue [<label>]
type ContinueStmt struct{ Kw, Label, Semi *lexing.Token }

// BreakStmt is the break statement
// break [<label>]
type BreakStmt struct{ Kw, Label, Semi *lexing.Token }

// FallthroughStmt is the fallthrough statement
// fallthrough
type FallthroughStmt struct{ Kw, Semi *lexing.Token }

// EmptyStmt is an empty statement created by
// an orphan semicolon
type EmptyStmt struct {
	Semi *lexing.Token
}
