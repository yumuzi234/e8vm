package ast

import (
	"fmt"

	"e8vm.io/e8vm/lex8"
)

// DeclPos returns the keyword position of a top-level declaration.
func DeclPos(d Decl) *lex8.Pos {
	switch d := d.(type) {
	case *VarDecls:
		return d.Kw.Pos
	case *ConstDecls:
		return d.Kw.Pos
	case *Func:
		return d.Kw.Pos
	case *Struct:
		return d.Kw.Pos
	}

	panic(fmt.Errorf("invalid top-level decl: %T", d))
}
