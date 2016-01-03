package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
)

type objVar struct {
	name string
	*ref
}

type objFunc struct {
	name string
	*ref
	f        *ast.Func
	isMethod bool
	isAlias  bool
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
