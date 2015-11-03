// Package ast declares the types for the abstract syntax tree in E8VM's
// assembly language.
package ast

import (
	"e8vm.io/e8vm/lex8"
)

// File represents a file.
type File struct {
	Imports *Import

	Decls    []interface{}
	Comments []*lex8.Token
}

// a listing of possible declarations
var decls = []interface{}{
	new(Func),
	new(Var),
}
