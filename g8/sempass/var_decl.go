package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

func buildVarDecl(b *Builder, d *ast.VarDecl) tast.Stmt {
	idents := d.Idents.Idents

	if d.Eq != nil {
		right := b.BuildExpr(d.Exprs)
		if right == nil {
			return nil
		}
		if d.Type != nil {
			tdest := b.BuildType(d.Type)
			if tdest != nil {
				return nil
			}
			if !types.IsAllocable(tdest) {
				pos := ast.ExprPos(d.Type)
				b.Errorf(pos, "%s is not allocatable", tdest)
				return nil
			}

			panic("todo")

		} else {
			return define(b, idents, right, d.Eq)
		}
	}

	if d.Type == nil {
		panic("type missing")
	}

	t := b.BuildType(d.Type)
	if t == nil {
		return nil
	}

	var syms []*sym8.Symbol
	for _, ident := range idents {
		s := declareVar(b, ident, t)
		if s == nil {
			return nil
		}
		syms = append(syms, s)
	}

	return &tast.DefineStmt{syms, nil}
}

func buildVarDecls(b *Builder, decls *ast.VarDecls) tast.Stmt {
	if len(decls.Decls) == 0 {
		return nil
	}

	var ret []tast.Stmt
	for _, d := range decls.Decls {
		d := buildVarDecl(b, d)
		if d == nil {
			return nil
		}
		ret = append(ret, d)
	}
	return &tast.VarDecls{ret}
}
