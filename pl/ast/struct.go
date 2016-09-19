package ast

import (
	"shanhu.io/smlvm/lexing"
)

// Field is a member variable of a struct
type Field struct {
	Idents *IdentList
	Type   Expr
	Semi   *lexing.Token
}

// Struct declares a structure type
type Struct struct {
	Kw      *lexing.Token
	Name    *lexing.Token
	KwAfter *lexing.Token
	Lbrace  *lexing.Token

	Fields  []*Field
	Methods []*Func

	Rbrace *lexing.Token
	Semi   *lexing.Token
}

// Interface is an interface type
// TODO:
/*
type Interface struct {
	Kw     *lex8.Token
	Name   *lex8.Token
	Lbrace *lex8.Token
	Funcs  []*FuncSig
	Rbrace *lex8.Token
}
*/
