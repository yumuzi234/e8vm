package link

import (
	"shanhu.io/smlvm/lexing"
)

// FuncMeta stores the meta data of a function
// for generating debug symbol.
type FuncMeta struct {
	Frame uint32
	Pos   *lexing.Pos
}
