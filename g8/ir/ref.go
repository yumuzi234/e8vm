package ir

// Ref is a general reference to a variable or constant
type Ref interface {
	Attr() *Attribute
}

// Attribute is a reference attribute
type Attribute struct {
	Size   int32
	Align  int32
	Signed bool
}

// Attr returns the attribute itself.
func (a *Attribute) Attr() *Attribute {
	return a
}

// Var is a local variable allocated on stack
type Var struct {
	Name string
	*Attribute
}

// AddrRef is an indirect reference of a variable
type AddrRef struct {
	Base   Ref
	Offset int32
	*Attribute
}

// Number is a constant number
type Number struct {
	V uint32
	*Attribute
}
