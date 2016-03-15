package g8

import (
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/codegen"
	"e8vm.io/e8vm/g8/types"
)

// ref is a reference to one or a list of typed objects.
type ref struct {
	lst []*ref

	typ         types.T
	addressable bool
	recv        *ref        // receiver, if any
	recvFunc    *types.Func // the actual func type

	ir codegen.Ref
}

func newRef(t types.T, r codegen.Ref) *ref { return &ref{typ: t, ir: r} }

func newTypeRef(t types.T) *ref { return &ref{typ: &types.Type{t}} }

func newAddressableRef(t types.T, r codegen.Ref) *ref {
	return &ref{typ: t, ir: r, addressable: true}
}

func newRecvRef(t *types.Func, recv *ref, r codegen.Ref) *ref {
	return &ref{typ: t.MethodFunc, ir: r, recv: recv, recvFunc: t}
}

func appendRef(r1, r2 *ref) *ref {
	if !r2.IsSingle() {
		panic("must merge single")
	}

	switch r1.Len() {
	case 0:
		return r2
	case 1:
		ref := new(ref)
		ref.lst = append(ref.lst, r1, r2)
		return ref
	default:
		r1.lst = append(r1.lst, r2)
		return r1
	}
}

func (r *ref) Len() int {
	if len(r.lst) == 0 {
		if r.typ == nil {
			return 0
		}
		return 1
	}

	return len(r.lst)
}

func (r *ref) IsSingle() bool { return r.typ != nil && len(r.lst) == 0 }
func (r *ref) IsConst() bool {
	return r.IsSingle() && types.IsConst(r.Type())
}

func (r *ref) Type() types.T {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.typ
}

func (r *ref) IsType() bool {
	if !r.IsSingle() {
		return false
	}
	_, ok := r.Type().(*types.Type)
	return ok
}

func (r *ref) IsPkg() bool {
	if !r.IsSingle() {
		return false
	}
	_, ok := r.Type().(*types.Pkg)
	return ok
}

func (r *ref) IsNil() bool {
	if !r.IsSingle() {
		return false
	}
	return types.IsNil(r.Type())
}

func (r *ref) TypeType() types.T {
	return r.Type().(*types.Type).T
}

func (r *ref) IR() codegen.Ref {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.ir
}

func (r *ref) Addressable() bool {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.addressable
}

func (r *ref) String() string {
	if r == nil {
		return "void"
	}

	if len(r.lst) == 0 {
		if r.typ == nil {
			return "<nil>"
		}
		return r.typ.String()
	}

	return fmt8.Join(r.TypeList(), ",")
}

func (r *ref) IsBool() bool {
	if !r.IsSingle() {
		return false
	}
	return types.IsBasic(r.Type(), types.Bool)
}

func (r *ref) At(i int) *ref {
	if r.IsSingle() {
		if i != 0 {
			panic("invalid index")
		}
		return r
	}

	return r.lst[i]
}

func (r *ref) IRList() []codegen.Ref {
	if len(r.lst) == 0 {
		if r.typ == nil {
			return nil
		}
		return []codegen.Ref{r.ir}
	}

	var ret []codegen.Ref
	for _, ref := range r.lst {
		ret = append(ret, ref.IR())
	}

	return ret
}

func (r *ref) TypeList() []types.T {
	if len(r.lst) == 0 {
		if r.typ == nil {
			return nil
		}
		return []types.T{r.typ}
	}

	var ret []types.T
	for _, ref := range r.lst {
		ret = append(ret, ref.Type())
	}
	return ret
}
