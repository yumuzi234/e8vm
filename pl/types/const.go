package types

import (
	"fmt"
)

// Const saves a compile time constant.
type Const struct {
	Type  T // optional type
	Value interface{}
}

// NewConstInt creates a new constant for a specific int type.
func NewConstInt(v int64, t T) *Const {
	if !IsInteger(t) {
		panic("the type for NewConstInt must be a int or unit")
	}
	return &Const{Value: v, Type: t}
}

// NewConstString creates a new string constant.
func NewConstString(s string) *Const {
	return &Const{Value: s, Type: String}
}

// NewConstBool creates a new bool constant.
func NewConstBool(v bool) *Const {
	return &Const{Value: v, Type: Bool}
}

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
	case string:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%s", v)
	}
}

// IsConst checks if a type is constant.
func IsConst(t T) bool {
	_, ok := t.(*Const)
	return ok
}

// ConstType checks and transforms a type to const type.
func ConstType(t T) (*Const, bool) {
	c, ok := t.(*Const)
	if !ok {
		return nil, false
	}
	return c, true
}

// NumConst checks and transforms a type to a typeless number.
// this function will be replaced by typed const
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
