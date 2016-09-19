package ast

import (
	"shanhu.io/smlvm/lexing"
)

// Func is an assembly function.
type Func struct {
	Stmts []*FuncStmt

	Kw, Name             *lexing.Token
	Lbrace, Rbrace, Semi *lexing.Token
}

// FuncStmt is a statement in a assembly function.
// It is either a instruction or a label.
type FuncStmt struct {
	Ops []*lexing.Token
}
