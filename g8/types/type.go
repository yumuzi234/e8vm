package types

import (
	"fmt"
)

// Type is a type type. It is the type of an expression in G language that
// represents a type. For example, "int" is an expressin of type type.
type Type struct {
	T // the type it represents.
}

// Size will panic.
func (*Type) Size() int32 { panic("bug") }

// RegSizeAlign will panic.
func (*Type) RegSizeAlign() bool { panic("bug") }

func (t *Type) String() string {
	return fmt.Sprintf("%s(type)", t.T.String())
}
