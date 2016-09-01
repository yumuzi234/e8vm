package ast

import (
	"path"
	"strconv"

	"e8vm.io/e8vm/lexing"
)

// ImportDecl is a import declare line
type ImportDecl struct {
	As   *lexing.Token // optional
	Path *lexing.Token
	Semi *lexing.Token
}

// ImportDecls is a top-level import declaration block
type ImportDecls struct {
	Kw     *lexing.Token
	Lparen *lexing.Token
	Decls  []*ImportDecl
	Rparen *lexing.Token
	Semi   *lexing.Token
}

// ImportPos returns the position of the import symbol
func ImportPos(d *ImportDecl) *lexing.Pos {
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
