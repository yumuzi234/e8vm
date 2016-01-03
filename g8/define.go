package g8

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func declareVar(b *builder, tok *lex8.Token, t types.T) *objVar {
	name := tok.Lit
	v := &objVar{name: name}
	s := sym8.Make(b.path, name, tast.SymVar, v, t, tok.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(tok.Pos, "%q already declared as a %s",
			name, tast.SymStr(conflict.Type),
		)
		return nil
	}
	return v
}

func declareVarRef(b *builder, tok *lex8.Token, r *ref) {
	obj := declareVar(b, tok, r.Type())
	if obj != nil {
		obj.ref = r
	}
}

func buildDefine(b *builder, d *tast.Define) {
	var refs []*ref
	for _, sym := range d.Left {
		name := sym.Name()
		t := sym.ObjType.(types.T)
		v := b.newLocal(t, name)
		r := newAddressableRef(t, v)
		sym.Obj = &objVar{name: name, ref: r}
		refs = append(refs, r)
	}

	n := len(refs)
	if d.Right == nil {
		for i := 0; i < n; i++ {
			b.b.Zero(refs[i].IR())
		}
	} else {
		src := b.buildExpr(d.Right)
		for i := 0; i < n; i++ {
			b.b.Assign(refs[i].IR(), src.At(i).IR())
		}
	}
}

func buildConstDefine(b *builder, d *tast.Define) {
	for _, sym := range d.Left {
		name := sym.Name()
		t := sym.ObjType.(types.T)
		r := newRef(t, nil)
		sym.Obj = &objConst{name: name, ref: r}
	}
}
