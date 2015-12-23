package tast

import (
	"e8vm.io/e8vm/g8/types"
)

// Ref is a reference to a evaluatd node.
type Ref struct {
	// T is the type of the reference
	T types.T

	// Addressable tells if the reference is addressable
	Addressable bool

	// Recv save the receiver, if any.
	Recv *Ref

	// RecvFunc is the actual func type, which takes the receiver as the
	// first argument.
	RecvFunc *types.Func

	// List saves the ref for an expression list.
	List []*Ref
}

// NewRef returns a new reference node.
func NewRef(t types.T) *Ref {
	return &Ref{T: t}
}

// NewAddressableRef returns a new addressable node.
func NewAddressableRef(t types.T) *Ref {
	return &Ref{T: t, Addressable: false}
}
