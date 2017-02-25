package types

import (
	"fmt"
	"math"
)

// Const saves a compile time constant.
type Const struct {
	Type  T // optional type
	Value interface{}
}

// NewConstInt creates a new constant for a specific int type.
func NewConstInt(v int64, t T) *Const {
	if !(IsInteger(t) && InRange(v, t)) {
		panic("the type for NewConstInt must be a int or unit")
	}
	return &Const{Value: v, Type: t}
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

// InRange checks if a const is in range of an integer type.
func InRange(v int64, t T) bool {
	t, ok := t.(Basic)
	if !ok {
		return false
	}
	switch t {
	case Int:
		return v >= math.MinInt32 && v <= math.MaxInt32
	case Uint:
		return v >= 0 && v <= math.MaxUint32
	case Int8:
		return v >= math.MinInt8 && v <= math.MaxInt8
	case Uint8:
		return v >= 0 && v <= math.MaxUint8
	}
	return false
}
