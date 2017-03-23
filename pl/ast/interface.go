package ast

import (
	"shanhu.io/smlvm/lexing"
)

// Interface declares a interface
type Interface struct {
	Kw       *lexing.Token
	Name     *lexing.Token
	Lbrace   *lexing.Token
	Funcs    []*lexing.Token
	FuncSigs []*FuncSig
	Rbrace   *lexing.Token
	Semi     *lexing.Token
}
