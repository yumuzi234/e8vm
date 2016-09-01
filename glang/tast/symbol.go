package tast

import (
	"fmt"
)

// Symbol types
const (
	SymNone = iota
	SymFunc
	SymVar
	SymStruct
	SymType
	SymConst
	SymImport
	SymField
)

// SymStr returns the string representation of a symbol
func SymStr(s int) string {
	switch s {
	case SymVar:
		return "variable"
	case SymStruct:
		return "struct"
	case SymFunc:
		return "function"
	case SymConst:
		return "constant"
	case SymImport:
		return "imported package"
	case SymField:
		return "struct field"
	case SymType:
		return "builtin type"
	default:
		panic(fmt.Errorf("unknown symbol: %d", s))
	}
}
