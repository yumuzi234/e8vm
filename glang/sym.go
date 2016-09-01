package glang

import (
	"e8vm.io/e8vm/glang/ast"
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

// TODO: this can be removed if typed const is fully supported in g8/types. We
// are keeping this now only for builtin consts like nil, true and false.
type objConst struct {
	name string
	*ref
}
