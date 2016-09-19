package codegen

import (
	"shanhu.io/smlvm/lexing"
)

type gener struct {
	memClear *FuncSym
	memCopy  *FuncSym

	*lexing.ErrorList
}

func newGener() *gener {
	return &gener{
		ErrorList: lexing.NewErrorList(),
	}
}
