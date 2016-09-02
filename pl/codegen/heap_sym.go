package codegen

// HeapSym is a variable on heap.
type HeapSym struct {
	pkg, name    string
	size         int32
	u8           bool
	regSizeAlign bool
}

// NewHeapSym creates a new heap var symbol.
func NewHeapSym(
	pkg, name string, size int32, u8, regSizeAlign bool,
) *HeapSym {
	return &HeapSym{
		pkg:          pkg,
		name:         name,
		size:         size,
		u8:           u8,
		regSizeAlign: regSizeAlign,
	}
}

func (s *HeapSym) String() string { return s.name }

// Size returns the size of the variable.
func (s *HeapSym) Size() int32 { return s.size }

// RegSizeAlign tells if the variable is word aligned.
func (s *HeapSym) RegSizeAlign() bool { return s.regSizeAlign }
