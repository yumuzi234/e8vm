package ast

import (
	"path"
	"strconv"

	"e8vm.io/e8vm/lex8"
)

// ImportDecl is a import declare line
type ImportDecl struct {
	As   *lex8.Token // optional
	Path *lex8.Token
	Semi *lex8.Token
}

// ImportDecls is a top-level import declaration block
type ImportDecls struct {
	Kw     *lex8.Token
	Lparen *lex8.Token
	Decls  []*ImportDecl
	Rparen *lex8.Token
	Semi   *lex8.Token
}

// ImportPos returns the position of the import symbol
func ImportPos(d *ImportDecl) *lex8.Pos {
	if d.As == nil {
		return d.Path.Pos
	}
	return d.As.Pos
}

// ImportPathAs parses the import path and as string.
func ImportPathAs(d *ImportDecl) (p, as string, err error) {
	p, err = strconv.Unquote(d.Path.Lit)
	if err != nil {
		return "", "", err
	}

	if d.As == nil {
		return p, path.Base(p), nil
	}
	return p, d.As.Lit, nil
}
