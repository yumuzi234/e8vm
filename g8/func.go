package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func makeRetRef(ts []*types.Arg, irs []ir.Ref) *ref {
	if len(ts) != len(irs) {
		panic("bug")
	}
	if len(ts) == 0 {
		return nil
	}

	ret := new(ref)
	for i, t := range ts {
		ret = appendRef(ret, newAddressableRef(t.T, irs[i]))
	}
	return ret
}

func buildFunc(b *builder, f *tast.Func, irFunc *ir.Func) {
	b.f = irFunc

	if f.Receiver != nil {
		// bind the receiver
		t := f.Receiver.ObjType.(types.T)
		ref := newAddressableRef(t, irFunc.ThisRef())
		f.Receiver.Obj = &objVar{f.Receiver.Name(), ref}
	} else if f.This != nil {
		// bind this pointer
		b.this = newRef(f.This, irFunc.ThisRef())
	}

	// bind arg symbols
	args := irFunc.ArgRefs()
	if f.IsMethod() {
		args = args[1:] // skip <this>
	}
	for i, s := range f.Args {
		if s != nil {
			t := s.ObjType.(types.T)
			ref := newAddressableRef(t, args[i])
			s.Obj = &objVar{s.Name(), ref}
		}
	}

	// bind named return symbols
	rets := irFunc.RetRefs()

	t := f.Sym.ObjType.(*types.Func)
	b.fretRef = makeRetRef(t.Rets, rets)
	if f.NamedRets != nil {
		for i, s := range f.NamedRets {
			if s != nil {
				s.Obj = &objVar{s.Name(), b.fretRef.At(i)}
			}
		}
	}

	b.b = b.f.NewBlock(nil)
	for _, stmt := range f.Body {
		b.buildStmt2(stmt)
	}
}
