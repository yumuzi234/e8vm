package ast

import (
	"e8vm.io/e8vm/lexing"
)

// Import is an import declaration block
type Import struct {
	Stmts []*ImportStmt

	Kw, Lbrace, Rbrace, Semi *lexing.Token
}

// ImportStmt is an import statement
type ImportStmt struct {
	Path *lexing.Token
	As   *lexing.Token
}
