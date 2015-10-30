package g8

import (
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

// ref is a reference to one or a list of typed objects.
type ref struct {
	typ         []types.T
	ir          []ir.Ref // this is essentially anything
	addressable []bool

	recv     *ref        // receiver, if any
	recvFunc *types.Func // the actual func type
}

func newSingleRef(t types.T, r ir.Ref, addressable bool) *ref {
	return &ref{
		typ:         []types.T{t},
		ir:          []ir.Ref{r},
		addressable: []bool{addressable},
	}
}

// newRef creates a simple single ref
func newRef(t types.T, r ir.Ref) *ref {
	return newSingleRef(t, r, false)
}

func newTypeRef(t types.T) *ref {
	return newRef(&types.Type{t}, nil)
}

func newAddressableRef(t types.T, r ir.Ref) *ref {
	return newSingleRef(t, r, true)
}

func newRecvRef(t *types.Func, recv *ref, r ir.Ref) *ref {
	ret := newRef(t.MethodFunc, r)
	ret.recv = recv
	ret.recvFunc = t
	return ret
}

func (r *ref) Len() int       { return len(r.typ) }
func (r *ref) IsSingle() bool { return len(r.typ) == 1 }
func (r *ref) IsConst() bool {
	return r.IsSingle() && types.IsConst(r.Type())
}

func (r *ref) Type() types.T {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.typ[0]
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

func (r *ref) IR() ir.Ref {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.ir[0]
}

func (r *ref) Addressable() bool {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.addressable[0]
}

func (r *ref) String() string {
	if r == nil {
		return "void"
	}

	if len(r.typ) == 0 {
		return "<nil>"
	}

	return fmt8.Join(r.typ, ",")
}

func (r *ref) IsBool() bool {
	if !r.IsSingle() {
		return false
	}
	return types.IsBasic(r.Type(), types.Bool)
}
