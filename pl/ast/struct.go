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

	Fields []*Field

	Rbrace *lexing.Token
	Semi   *lexing.Token
}
