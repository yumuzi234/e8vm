package ast

import (
	"e8vm.io/e8vm/lexing"
)

// Var is a variable declaration
type Var struct {
	Stmts []*VarStmt

	Kw, Name             *lexing.Token
	Lbrace, Rbrace, Semi *lexing.Token
}

// VarStmt is a variable statement.
type VarStmt struct {
	Type *lexing.Token
	Args []*lexing.Token
}
