package types

import (
	"fmt"
)

// Basic types are the built-in fundamentals of the language
type Basic int

// Basic type codes
const (
	Int Basic = iota
	Uint
	Int8
	Uint8
	Float32
	Bool
)

// IsBasic checks if a type is a particular basic type
func IsBasic(t T, b Basic) bool {
	code, ok := t.(Basic)
	if !ok {
		return false
	}
	return code == b
}

// IsRegSizeBasic checks if a type is a Int, Uint or Float32
func IsRegSizeBasic(t T) bool {
	code, ok := t.(Basic)
	if !ok {
		return false
	}
	switch code {
	case Int, Uint, Float32:
		return true
	}
	return false
}

// IsInteger checks if a type is an integer type
func IsInteger(t T) bool {
	code, ok := t.(Basic)
	if !ok {
		return false
	}
	switch code {
	case Int, Uint, Int8, Uint8:
		return true
	}
	return false
}

// IsSigned checks if a type is a signed integer type
func IsSigned(t T) bool {
	code, ok := t.(Basic)
	if !ok {
		return false
	}
	switch code {
	case Int, Int8:
		return true
	}
	return false
}

// IsUnsigned checks if a type is an unsigned integer type
func IsUnsigned(t T) bool {
	code, ok := t.(Basic)
	if !ok {
		return false
	}
	switch code {
	case Uint, Uint8:
		return true
	}
	return false
}

// IsByte checks if a type is Uint8
func IsByte(t T) bool {
	code, ok := t.(Basic)
	if !ok {
		return false
	}
	return code == Uint8
}

// BothBasic checks if two types are both a particular basic type
func BothBasic(a, b T, t Basic) bool {
	return IsBasic(a, t) && IsBasic(b, t)
}

// SameBasic check if two types are both basic types, and also returns
// the type if it is.
func SameBasic(a, b T) (bool, Basic) {
	code1, ok := a.(Basic)
	if !ok {
		return false, Int
	}
	code2, ok := b.(Basic)
	if !ok {
		return false, Int
	}
	if code1 != code2 {
		return false, Int
	}
	return true, code1
}

// Size returns the size in memory of a basic type
func (t Basic) Size() int32 {
	switch t {
	case Int, Uint:
		return 4
	case Int8, Uint8:
		return 1
	case Float32:
		return 4
	case Bool:
		return 1
	default:
		panic("unknown basic type")
	}
}

// String returns the name of the basic type.
func (t Basic) String() string {
	switch t {
	case Int:
		return "int"
	case Uint:
		return "uint"
	case Int8:
		return "int8"
	case Uint8:
		return "uint8"
	case Bool:
		return "bool"
	case Float32:
		return "float32"
	default:
		panic(fmt.Errorf("invalid basic type %d", t))
	}
}

// RegSizeAlign checks if the type is word aligned.
func (t Basic) RegSizeAlign() bool {
	switch t {
	case Int, Uint, Float32:
		return true
	case Int8, Uint8, Bool:
		return false
	default:
		panic(fmt.Errorf("invalid basic type %d", t))
	}
}
