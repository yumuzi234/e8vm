package debug8

import (
	"e8vm.io/e8vm/lex8"
)

// Func saves the debug information of a function
type Func struct {
	Start uint32
	Size  uint32
	Frame uint32

	Pos *lex8.Pos
}
