// Package tast saves data structures for typed AST.
// Compiling it should contain no semantic errors.
package tast

import (
	"shanhu.io/smlvm/pl/types"
)

// Expr is a generic interface for an expression.
type Expr interface {
	R() *Ref
	Type() types.T
}

// Stmt is a generic interface for a statement.
type Stmt interface{}
