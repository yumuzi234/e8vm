// Package tast saves data structures for typed AST.
// Compiling it should contain no semantic errors.
package tast

// Expr is a generic interface for an expression.
type Expr interface {
	R() *Ref
}

// Stmt is a generic interface for a statement.
type Stmt interface{}
