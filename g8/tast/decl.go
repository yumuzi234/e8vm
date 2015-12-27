package tast

import (
	"e8vm.io/e8vm/g8/types"
)

// VarDecl declares a set of variables.
type VarDecl struct {
	Idents []*Ref
	Type   types.T
	Exprs  *ExprList
}

// VarDecls is a variable declaration group.
type VarDecls struct {
	Decls []Stmt
}

// ConstDecl declares a set of constants.
type ConstDecl struct {
	Idents []*Ref
	Type   types.T
	Exprs  *ExprList
}

// ConstDecls is a const declaration group.
type ConstDecls struct {
	Decls []*ConstDecl
}
