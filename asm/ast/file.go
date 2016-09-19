// Package ast declares the types for the abstract syntax tree in E8VM's
// assembly language.
package ast

import (
	"shanhu.io/smlvm/lexing"
)

// File represents a file.
type File struct {
	Imports *Import

	Decls    []interface{}
	Comments []*lexing.Token
}

// a listing of possible declarations
var decls = []interface{}{
	new(Func),
	new(Var),
}
