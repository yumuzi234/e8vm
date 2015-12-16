package ir

// FuncPtr is a reference to a function pointer.
type FuncPtr struct {
	sig *FuncSig
	Ref
}

// NewFuncPtr wraps a new function pointer.
func NewFuncPtr(sig *FuncSig, ref Ref) *FuncPtr {
	return &FuncPtr{sig, ref}
}
