package ast

import (
	"e8vm.io/e8vm/lexing"
)

// FuncRecv is the receiver of a struct method
type FuncRecv struct {
	Lparen     *lexing.Token
	Recv       *lexing.Token
	Star       *lexing.Token
	StructName *lexing.Token
	Rparen     *lexing.Token
}

// FuncAlias is for aliasing an imported function
type FuncAlias struct {
	Eq   *lexing.Token
	Pkg  *lexing.Token
	Dot  *lexing.Token
	Name *lexing.Token
}

// Func is a function
type Func struct {
	Kw   *lexing.Token
	Name *lexing.Token

	Recv *FuncRecv
	*FuncSig

	Alias *FuncAlias
	Body  *Block
	Semi  *lexing.Token
}
