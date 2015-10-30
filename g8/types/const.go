package types

import (
	"fmt"
)

// Const saves a compile time constant.
type Const struct {
	Type  T // optional type
	Value interface{}
}

// NewConst creates a new constant.
func NewConst(v int64, t T) *Const { return &Const{Value: v, Type: t} }

// NewNumber creates a new constant number.
func NewNumber(v int64) *Const {
	return &Const{Value: v, Type: Number{}}
}

// Size returns the type of the size.
func (c *Const) Size() int32 { return c.Type.Size() }

// RegSizeAlign is a shortcut for c.T.RegSizeAlign()
func (c *Const) RegSizeAlign() bool { return c.Type.RegSizeAlign() }

// String returns the number
func (c *Const) String() string {
	switch v := c.Value.(type) {
	case int64:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%s", v)
	}
}

// IsConst checks if a type is constant.
func IsConst(t T) bool {
	_, ok := t.(*Const)
	return ok
}

// NumConst checks and transforms a type to a typeless number.
func NumConst(t T) (int64, bool) {
	c, ok := t.(*Const)
	if !ok {
		return 0, false
	}
	if _, ok := c.Type.(Number); !ok {
		return 0, false
	}
	return c.Value.(int64), true
}
