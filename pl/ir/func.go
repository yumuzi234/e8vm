package ir

import (
	"shanhu.io/smlvm/arch"
)

// FuncSig is a function signature
type FuncSig struct {
	Args []*Var
	Rets []*Var
}

// Func is a function
type Func struct {
	Name string

	*FuncSig
	Locals []*Var
	Blocks []*Block
}

// FuncSym is a function symbol
type FuncSym struct {
	*Symbol
}

var funcAttr = &Attribute{
	Size:  arch.RegSize,
	Align: arch.RegSize,
}

// Attr returns the attribute of a function symbol
func (s *FuncSym) Attr() *Attribute {
	return funcAttr
}
