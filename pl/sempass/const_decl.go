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

	var ret []*syms.Symbol
	var tdest types.T
	if d.Type != nil {
		tdest = b.buildType(d.Type)
		if tdest == nil {
			return nil
		}
		if _, ok := tdest.(types.Basic); !ok {
			b.CodeErrorf(ast.ExprPos(d.Type), "pl.constType",
				"%s is not supported for a const type", tdest)
			return nil
		}
	}
	var right tast.Expr
	idents := d.Idents.Idents

	if d.Eq != nil {
		right = buildConstExprList(b, d.Exprs)
		if right == nil {
			return nil
		}

		nright := right.R().Len()
		nleft := len(idents)
		if nleft != nright {
			b.CodeErrorf(d.Eq.Pos, "pl.cannotAssign.lengthMismatch",
				"cannot assign(len) %s to %s; length mismatch",
				nright, nleft)
			return nil
		}
	}

	for i, ident := range idents {
		var sym *syms.Symbol
		var t types.T
		if right == nil {
			t, _ = types.NewConstInt(0, tdest)
		} else {
			t = right.R().At(i).Type()
			if !types.IsConst(t) {
				b.CodeErrorf(ast.ExprPos(d.Exprs.Exprs[i]),
					"pl.expectConstExpr", "not a const")
				return nil
			}
			ct, _ := t.(*types.Const)
			if tdest != nil {
				t = types.CastConst(ct, tdest)
			}
		}

		sym = declareConst(b, ident, t)
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
