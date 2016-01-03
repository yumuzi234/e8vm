package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func declareVar(b *Builder, tok *lex8.Token, t types.T) *sym8.Symbol {
	name := tok.Lit
	s := sym8.Make(b.path, name, tast.SymVar, nil, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already defined as a %s",
			name, tast.SymStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously defined here")
		return nil
	}
	return s
}

func declareVars(b *Builder, ids []*lex8.Token, t types.T) []*sym8.Symbol {
	var syms []*sym8.Symbol
	for _, id := range ids {
		s := declareVar(b, id, t)
		if s == nil {
			return nil
		}
		syms = append(syms, s)
	}
	return syms
}
