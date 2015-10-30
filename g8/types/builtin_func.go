package types

import (
	"fmt"
)

// BuiltInFunc is the type of a builtin function
type BuiltInFunc struct {
	Name string
}

// NewBuiltInFunc creates the type of a particular builtin function
func NewBuiltInFunc(name string) *BuiltInFunc {
	ret := new(BuiltInFunc)
	ret.Name = name
	return ret
}

// Size on a builtin function type will panic.
func (f *BuiltInFunc) Size() int32 { panic("size on builtin func") }

func (f *BuiltInFunc) String() string {
	return fmt.Sprintf("%s() (builtin)", f.Name)
}

// RegSizeAlign on a builtin function will panic.
func (f *BuiltInFunc) RegSizeAlign() bool {
	panic("reg align on builtin func")
}
