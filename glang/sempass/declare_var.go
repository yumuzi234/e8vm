package sempass

import (
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/syms"
)

func declareVar(
	b *builder, tok *lexing.Token, t types.T, used bool,
) *syms.Symbol {
	name := tok.Lit
	s := syms.Make(b.path, name, tast.SymVar, nil, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already defined as a %s",
			name, tast.SymStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously defined here")
		return nil
	}
	s.Used = used
	return s
}

func declareVars(
	b *builder, ids []*lexing.Token, t types.T, used bool,
) []*syms.Symbol {
	var ret []*syms.Symbol
	for _, id := range ids {
		s := declareVar(b, id, t, used)
		if s == nil {
			return nil
		}
		ret = append(ret, s)
	}
	return ret
}
