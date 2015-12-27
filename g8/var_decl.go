package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

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
