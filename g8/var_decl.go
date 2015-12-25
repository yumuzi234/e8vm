package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func allocTypedVars(b *builder, toks []*lex8.Token, t types.T) *ref {
	ts := make([]types.T, len(toks))
	for i := range toks {
		ts[i] = t
	}
	return allocVars(b, toks, ts)
}

func zero(b *builder, ref *ref) {
	for _, r := range ref.IRList() {
		b.b.Zero(r)
	}
}

func buildVarDecl(b *builder, d *ast.VarDecl) {
	idents := d.Idents.Idents

	if d.Eq != nil {
		right := buildExprList(b, d.Exprs)
		if right == nil {
			return
		}
		if d.Type != nil {
			tdest := b.spass.BuildType(d.Type)
			if tdest == nil {
				return
			}

			dest := allocTypedVars(b, idents, tdest)
			if dest == nil {
				return
			}

			if assign(b, dest, right, d.Eq) {
				declareVars(b, idents, dest)
			}
		} else {
			define(b, idents, right, d.Eq)
		}
		return
	}

	if d.Type == nil {
		panic("type missing")
	}

	t := b.spass.BuildType(d.Type)
	if t == nil {
		return
	}

	for _, ident := range idents {
		r := newAddressableRef(t, b.newLocal(t, ident.Lit))
		zero(b, r)
		declareVarRef(b, ident, r)
	}
}

func buildVarDecls(b *builder, decls *ast.VarDecls) {
	for _, d := range decls.Decls {
		buildVarDecl(b, d)
	}
}

func buildGlobalVarDecl(b *builder, d *ast.VarDecl) {
	if d.Eq != nil {
		b.Errorf(d.Eq.Pos, "init for global var not supported yet")
		return
	}
	if d.Type == nil {
		panic("type missing")
	}

	t := b.spass.BuildType(d.Type)
	if t == nil {
		return
	}

	for _, ident := range d.Idents.Idents {
		obj := declareVar(b, ident, t)
		if obj != nil {
			obj.ref = newAddressableRef(t, b.newGlobalVar(t, ident.Lit))
		}
	}
}
