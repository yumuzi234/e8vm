package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

func declareConst(b *builder, tok *lexing.Token, t types.T) *syms.Symbol {
	name := tok.Lit
	s := syms.Make(b.path, name, tast.SymConst, nil, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.CodeErrorf(tok.Pos, "pl.declConflict.const",
			"%q already declared as a %s", name, tast.SymStr(conflict.Type),
		)
		b.CodeErrorf(conflict.Pos, "pl.declConflict.previousPos",
			"previously defined here")
		return nil
	}
	return s
}

func buildConstDecl(b *builder, d *ast.ConstDecl) *tast.Define {
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
	var tdest types.T
	if d.Type != nil {
		tdest = b.buildType(d.Type)
		if tdest == nil {
			return nil
		}
		if _, ok := tdest.(*types.Basic); !ok {
			b.CodeErrorf(ast.ExprPos(d.Type), "pl.constType",
				"%s is not support for a const type", tdest)
			return nil
		}
	}

	for i, ident := range idents {
		t := right.R().At(i).Type()
		if !types.IsConst(t) {
			b.CodeErrorf(ast.ExprPos(d.Exprs.Exprs[i]), "pl.expectConstExpr",
				"not a const")
			return nil
		}
		ct, _ := t.(*types.Const)
		if d.Type != nil {
			if !types.DefineTypeConst(tdest, ct) {
				b.CodeErrorf(ast.ExprPos(d.Exprs.Exprs[i]),
					"pl.cannotAssign.typeMismatch",
					"cannot assign const %s with %s", tdest, ct.Type)
				return nil
			}
			t = tdest
		}
		sym := declareConst(b, ident, t)
		if sym == nil {
			return nil
		}
		ret = append(ret, sym)
	}

	return &tast.Define{Left: ret, Right: right}
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
	return &tast.ConstDecls{Decls: ret}
}
