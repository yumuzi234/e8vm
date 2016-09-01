package codegen

import (
	"e8vm.io/e8vm/lexing"
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
