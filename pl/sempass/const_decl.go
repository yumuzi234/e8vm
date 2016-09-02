package sempass

import (
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/pl/ast"
	"e8vm.io/e8vm/pl/tast"
	"e8vm.io/e8vm/pl/types"
	"e8vm.io/e8vm/syms"
)

func declareConst(b *builder, tok *lexing.Token, t types.T) *syms.Symbol {
	name := tok.Lit
	s := syms.Make(b.path, name, tast.SymConst, nil, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already declared as a %s",
			name, tast.SymStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously defined here")
		return nil
	}
	return s
}

func buildConstDecl(b *builder, d *ast.ConstDecl) *tast.Define {
	if d.Type != nil {
		b.Errorf(ast.ExprPos(d.Type), "typed const not implemented yet")
		return nil
	}

	right := buildConstExprList(b, d.Exprs)
	if right == nil {
		return nil
	}

	nright := right.R().Len()
	idents := d.Idents.Idents
	nleft := len(idents)
	if nleft != nright {
		b.Errorf(d.Eq.Pos, "%d values for %d identifiers",
			nright, nleft,
		)
		return nil
	}

	var ret []*syms.Symbol
	for i, ident := range idents {
		t := right.R().At(i).Type()
		if !types.IsConst(t) {
			b.Errorf(ast.ExprPos(d.Exprs.Exprs[i]), "not a const")
			return nil
		}

		sym := declareConst(b, ident, t)
		if sym == nil {
			return nil
		}
		ret = append(ret, sym)
	}

	return &tast.Define{ret, right}
}

func buildConstDecls(b *builder, decls *ast.ConstDecls) tast.Stmt {
	if len(decls.Decls) == 0 {
		return nil
	}

	var ret []*tast.Define
	for _, d := range decls.Decls {
		d := buildConstDecl(b, d)
		if d != nil {
			ret = append(ret, d)
		}
	}
	return &tast.ConstDecls{ret}
}
