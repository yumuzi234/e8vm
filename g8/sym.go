package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
)

const (
	symNone = iota
	symFunc
	symVar
	symStruct
	symType
	symConst
	symImport
	symField
)

func symStr(s int) string {
	switch s {
	case symVar:
		return "variable"
	case symStruct:
		return "struct"
	case symFunc:
		return "function"
	case symConst:
		return "constant"
	case symImport:
		return "imported package"
	case symField:
		return "struct field"
	case symType:
		return "builtin type"
	default:
		panic(fmt.Errorf("unknown symbol: %d", s))
	}
}

type objVar struct {
	name string
	*ref
}

type objFunc struct {
	name string
	*ref
	f        *ast.Func
	isMethod bool
}

type objField struct {
	name string
	*types.Field
}

type objConst struct {
	name string
	*ref
}

type objType struct {
	name string
	*ref
}

type objImport struct {
	*ref
}
