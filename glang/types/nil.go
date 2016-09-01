package types

type null struct{}

var theNil T = null{}

// Nil returns the one and only nil
func Nil() T { return theNil }

// IsNil checks if it is the one and only nil
func IsNil(t T) bool { return t == theNil }

func (null) Size() int32 { panic("size on nil") }

func (null) String() string { return "nil" }

func (null) RegSizeAlign() bool { return true }
