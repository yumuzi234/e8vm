package ast

import (
	"shanhu.io/smlvm/lexing"
)

// Interface declares a interface
type Interface struct {
	Kw     *lexing.Token
	Name   *lexing.Token
	Lbrace *lexing.Token
	Funcs  []*InterfaceFunc
	Rbrace *lexing.Token
	Semi   *lexing.Token
}

// InterfaceFunc is a func in interface
type InterfaceFunc struct {
	Name     *lexing.Token
	FuncSigs *FuncSig
	Semi     *lexing.Token
}
