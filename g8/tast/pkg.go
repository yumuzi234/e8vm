package tast

import (
	"e8vm.io/e8vm/sym8"
)

// Import is an import statement
type Import struct {
	Sym *sym8.Symbol
}

// Pkg is a package of imports, consts, structs, vars and funcs.
type Pkg struct {
	Vars    []*VarDecls
	Funcs   []*Func
	Methods []*Func
}
