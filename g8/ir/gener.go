package ir

import (
	"e8vm.io/e8vm/lex8"
)

type gener struct {
	memClear *FuncSym
	memCopy  *FuncSym

	*lex8.ErrorList
}

func newGener() *gener {
	return &gener{
		ErrorList: lex8.NewErrorList(),
	}
}
