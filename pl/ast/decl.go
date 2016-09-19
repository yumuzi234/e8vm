package ast

import (
	"shanhu.io/smlvm/lexing"
)

// IdentList is a list of identifiers
type IdentList struct {
	Idents []*lexing.Token
	Commas []*lexing.Token
}

// VarDecl declares a set of variables. It is both a top level declaration
// and a statement.
type VarDecl struct {
	Idents *IdentList
	Type   Expr
	Eq     *lexing.Token
	Exprs  *ExprList
	Semi   *lexing.Token
}

// VarDecls is a variable declaration with a leading keyword
// It could be a single decl or a decl block.
type VarDecls struct {
	Kw     *lexing.Token
	Lparen *lexing.Token // optional
	Decls  []*VarDecl
	Rparen *lexing.Token // optional
	Semi   *lexing.Token
}

// ConstDecl declares a set of constants.
type ConstDecl struct {
	Idents *IdentList
	Type   Expr
	Eq     *lexing.Token
	Exprs  *ExprList
	Semi   *lexing.Token
}

// ConstDecls is a const declaration with a leading keyword
// It could be a single decl or a decl block.
type ConstDecls struct {
	Kw     *lexing.Token
	Lparen *lexing.Token // optional
	Decls  []*ConstDecl
	Rparen *lexing.Token // optional
	Semi   *lexing.Token
}
