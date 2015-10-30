package ir

import (
	"fmt"
)

type number struct{ v uint32 } // a constant number

func (n *number) String() string     { return fmt.Sprintf("%d", n.v) }
func (n *number) Size() int32        { return 4 }
func (n *number) RegSizeAlign() bool { return true }

// Num creates a constant reference to a int32 number
func Num(v uint32) Ref { return &number{v} }

// Snum creates a constant reference to a uint32 number
func Snum(v int32) Ref { return &number{uint32(v)} }

// SnumValue returns true, and the value of the number if it is a number.
// It returns false and 0 if it is not.
func SnumValue(r Ref) (bool, int32) {
	n, ok := r.(*number)
	if !ok {
		return false, 0
	}
	return true, int32(n.v)
}
