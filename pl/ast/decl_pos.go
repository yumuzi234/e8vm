package ast

import (
	"fmt"

	"shanhu.io/smlvm/lexing"
)

// DeclPos returns the keyword position of a top-level declaration.
func DeclPos(d Decl) *lexing.Pos {
	switch d := d.(type) {
	case *VarDecls:
		return d.Kw.Pos
	case *ConstDecls:
		return d.Kw.Pos
	case *Func:
		return d.Kw.Pos
	case *Struct:
		return d.Kw.Pos
	case *Interface:
		return d.Kw.Pos
	default:
		panic(fmt.Errorf("invalid top-level declaration type: %T", d))
	}
}
