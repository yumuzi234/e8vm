package ast

import (
	"e8vm.io/e8vm/lex8"
)

// Field is a member variable of a struct
type Field struct {
	Idents *IdentList
	Type   Expr
	Semi   *lex8.Token
}

// Struct declares a structure type
type Struct struct {
	Kw      *lex8.Token
	Name    *lex8.Token
	KwAfter *lex8.Token
	Lbrace  *lex8.Token

	Fields  []*Field
	Methods []*Func

	Rbrace *lex8.Token
	Semi   *lex8.Token
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
