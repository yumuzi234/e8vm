package tast

import (
	"shanhu.io/smlvm/syms"
)

// Import is an import statement
type Import struct {
	Sym *syms.Symbol
}

// Pkg is a package of imports, consts, structs, vars and funcs.
type Pkg struct {
	Imports   []*syms.Symbol
	Consts    []*syms.Symbol
	Structs   []*syms.Symbol
	Interface []*syms.Symbol

	Vars        []*Define
	FuncAliases []*FuncAlias
	Funcs       []*Func
	Methods     []*Func
}
