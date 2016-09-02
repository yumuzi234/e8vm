package codegen

import (
	"fmt"
)

// NewFuncSym creates a function symbol reference to a linkable function.
// It is used to perform function call operations to functions
// from other packages (functinos not declared in the current package,
// and hence only has a symbol and function signature).
func NewFuncSym(pkg, sym string, sig *FuncSig) Ref {
	return &FuncSym{pkg, sym, sig}
}

// FuncSym is a a function symbol
type FuncSym struct {
	pkg  string
	name string
	sig  *FuncSig
}

func (s *FuncSym) String() string {
	return fmt.Sprintf("%s.%s", s.pkg, s.name)
}

// Size returns the size of a function pointer.
func (s *FuncSym) Size() int32 { return regSize }

// RegSizeAlign returns true.
func (s *FuncSym) RegSizeAlign() bool { return true }
