package ast

import (
	"e8vm.io/e8vm/lex8"
)

// FuncRecv is the receiver of a struct method
type FuncRecv struct {
	Lparen     *lex8.Token
	Recv       *lex8.Token
	Star       *lex8.Token
	StructName *lex8.Token
	Rparen     *lex8.Token
}

// FuncAlias is for aliasing an imported function
type FuncAlias struct {
	Eq   *lex8.Token
	Pkg  *lex8.Token
	Dot  *lex8.Token
	Name *lex8.Token
}

// Func is a function
type Func struct {
	Kw   *lex8.Token
	Name *lex8.Token

	Recv *FuncRecv
	*FuncSig

	Alias *FuncAlias
	Body  *Block
	Semi  *lex8.Token
}
