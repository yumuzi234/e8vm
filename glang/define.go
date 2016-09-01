package glang

import (
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
)

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
