package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildVarDecl(b *builder, d *ast.VarDecl) *tast.Define {
	ids := d.Idents.Idents

	if d.Eq != nil {
		right := b.buildExpr(d.Exprs)
		if right == nil {
			return nil
		}

		if d.Type == nil {
			ret := define(b, ids, right, d.Eq)
			if ret == nil {
				return nil
			}
			return ret
		}

		tdest := b.buildType(d.Type)
		if tdest == nil {
			return nil
		}

		if !types.IsAllocable(tdest) {
			pos := ast.ExprPos(d.Type)
			b.CodeErrorf(pos, "pl.cannotAlloc",
				"type %s is not allocatable", tdest)
			return nil
		}

		// assignable check
		ts := right.R().TypeList()
		res := canAssignType(
			b, d.Eq.Pos, tdest, ts, "variable declaration",
		)
		if res.err {
			return nil
		}
		if res.needCast {
			right = tast.NewMultiCastType(right, tdest, res.castMask)
		}

		syms := declareVars(b, ids, tdest, false)
		if syms == nil {
			return nil
		}
		return &tast.Define{Left: syms, Right: right}
	}

	if d.Type == nil {
		panic("type missing")
	}

	t := b.buildType(d.Type)
	if t == nil {
		return nil
	}

	syms := declareVars(b, ids, t, false)
	if syms == nil {
		return nil
	}
	return &tast.Define{Left: syms, Right: nil}
}

func buildVarDecls(b *builder, decls *ast.VarDecls) tast.Stmt {
	if len(decls.Decls) == 0 {
		return nil
	}

	var ret []*tast.Define
	for _, d := range decls.Decls {
		d := buildVarDecl(b, d)
		if d != nil {
			ret = append(ret, d)
		}
	}
	return &tast.VarDecls{Decls: ret}
}
