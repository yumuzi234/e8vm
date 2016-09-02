package sempass

import (
	"e8vm.io/e8vm/pl/ast"
	"e8vm.io/e8vm/pl/tast"
)

func buildPkgVars(b *builder, vars []*ast.VarDecls) []*tast.Define {
	var ret []*tast.Define
	for _, decls := range vars {
		for _, d := range decls.Decls {
			v := buildPkgVar(b, d)
			if v != nil {
				ret = append(ret, v)
			}
		}
	}
	return ret
}

func buildPkgVar(b *builder, d *ast.VarDecl) *tast.Define {
	if d.Eq != nil {
		b.Errorf(d.Eq.Pos, "init for global var not supported yet")
		return nil
	}

	// since there's no init, we cannot possibly infer the type
	// from the expression
	if d.Type == nil {
		panic("type missing")
	}

	t := b.buildType(d.Type)
	if t == nil {
		return nil
	}

	syms := declareVars(b, d.Idents.Idents, t, false)
	if syms == nil {
		return nil
	}
	return &tast.Define{syms, nil}
}
