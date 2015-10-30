package ir

// HeapSym is a variable on heap.
type HeapSym struct {
	pkg, sym uint32 // base

	name         string
	size         int32
	u8           bool
	regSizeAlign bool
}

// NewHeapSym creates a new heap var symbol.
func NewHeapSym(
	name string, pkg, sym uint32, size int32, u8, regSizeAlign bool,
) *HeapSym {
	return &HeapSym{
		name:         name,
		pkg:          pkg,
		sym:          sym,
		size:         size,
		u8:           u8,
		regSizeAlign: regSizeAlign,
	}
}

func newHeapSym(size int32, name string, u8, regSizeAlign bool) *HeapSym {
	return &HeapSym{
		name: name, size: size, u8: u8,
		regSizeAlign: regSizeAlign,
	}
}

func (s *HeapSym) String() string { return s.name }

// Size returns the size of the variable.
func (s *HeapSym) Size() int32 { return s.size }

// RegSizeAlign tells if the variable is word aligned.
func (s *HeapSym) RegSizeAlign() bool { return s.regSizeAlign }

// Import returns a copy of the heap sym which package index is the given
// pindex.
func (s *HeapSym) Import(pindex uint32) *HeapSym {
	ret := *s
	ret.pkg = pindex
	return &ret
}
