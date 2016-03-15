package ir

// Ref is a general reference to a variable or constant
type Ref interface{}

// Attr is a reference attribute
type Attr struct {
	Size   int32
	Align  int32
	Signed bool
}

// Var is a local variable allocated on stack
type Var struct {
	Name string
	*Attr
}

// AddrRef is an indirect reference of a variable
type AddrRef struct {
	Base   Ref
	Offset int32
	*Attr
}

// Number is a constant number
type Number struct {
	V uint32
}
