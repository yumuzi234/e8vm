package tast

import (
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/glang/types"
)

// Ref is a reference to a evaluatd node.
type Ref struct {
	// T is the type of the reference
	T types.T

	// ConstValue is not nil when the reference already has a value.
	ConstValue interface{}

	// Addressable tells if the reference is addressable
	Addressable bool

	// Recv save the receiver, if any.
	Recv *Ref

	// List saves the ref for an expression list.
	List []*Ref

	// Void is true when the reference is a void ref.
	// This is required for saving ref of an expression list
	// or a return value of a function call.
	Void bool
}

// NewRef creates a new reference node.
func NewRef(t types.T) *Ref { return &Ref{T: t} }

// NewTypeRef creates a new reference node for a type expression.
func NewTypeRef(t types.T) *Ref { return NewRef(&types.Type{t}) }

// NewConstRef creates a new reference node with a constant value.
func NewConstRef(t types.T, v interface{}) *Ref {
	return &Ref{T: t, ConstValue: v}
}

// NewAddressableRef creates a new addressable node.
func NewAddressableRef(t types.T) *Ref {
	return &Ref{T: t, Addressable: true}
}

// Void is a void ref.
var Void = &Ref{Void: true}

// R returns itself.
func (r *Ref) R() *Ref { return r }

// At returns the ref in the list.
func (r *Ref) At(i int) *Ref {
	n := r.Len()
	if i < 0 || i >= n {
		panic("overflow")
	}
	if n == 1 {
		return r
	}

	return r.List[i]
}

// Type returns the type of the ref.
func (r *Ref) Type() types.T {
	if !r.IsSingle() {
		panic("not single")
	}
	return r.T
}

func (r *Ref) String() string {
	if r == nil {
		return "<error>"
	}
	if r.Void {
		return "void"
	}
	if len(r.List) == 0 {
		return r.T.String()
	}
	return fmt8.Join(r.List, ",")
}

// Len returns the number of refs in this ref bundle.
func (r *Ref) Len() int {
	if r.Void {
		return 0
	}
	if r.List != nil {
		return len(r.List)
	}
	return 1
}

// IsSingle checks if the ref is a single ref.
func (r *Ref) IsSingle() bool { return r.Len() == 1 }

// IsBool checks if the ref is a single boolean ref.
func (r *Ref) IsBool() bool {
	return r.IsSingle() && types.IsBasic(r.Type(), types.Bool)
}

// TypeList returns the type list of the reference's type.
func (r *Ref) TypeList() []types.T {
	if r.Void {
		return nil
	}
	if r.IsSingle() {
		return []types.T{r.T}
	}
	var ret []types.T
	for _, r := range r.List {
		ret = append(ret, r.Type())
	}
	return ret
}

// AppendRef append a ref in a ref bundle.
func AppendRef(base, toAdd *Ref) *Ref {
	if base.Void {
		return toAdd
	}
	if base.List == nil {
		return &Ref{
			List: []*Ref{base, toAdd},
		}
	}

	base.List = append(base.List, toAdd)
	return base
}
