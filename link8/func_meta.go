package link8

import (
	"e8vm.io/e8vm/lex8"
)

// FuncMeta stores the meta data of a function
// for generating debug symbol.
type FuncMeta struct {
	Frame uint32
	Pos   *lex8.Pos
}
