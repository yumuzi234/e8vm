package types

// Number is a typeless constant number.
type Number struct{}

// Size on Number will panic.
func (n Number) Size() int32 { panic("size of typeless number") }

// RegSizeAlign on Number will panic.
func (n Number) RegSizeAlign() bool { panic("alignment on constant") }

// String returns "<number>"
func (n Number) String() string { return "<number>" }
