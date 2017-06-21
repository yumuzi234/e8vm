package link

import (
	"math"
)

// Func is a relocatable code section
type Func struct {
	insts []uint32
	links []*link

	// filled when linking
	// TODO: this should not be here.
	addr uint32
}

// NewFunc creates a new relocatable code section.
func NewFunc() *Func {
	return new(Func)
}

// AddInst appends an instruction at the end of the function.
func (f *Func) AddInst(i uint32) {
	f.insts = append(f.insts, i)
}

// TooLarge checks if the function size is larger than 4GB.
func (f *Func) TooLarge() bool {
	return len(f.insts)*4 >= math.MaxInt32
}

// Size returns the size of the function.
func (f *Func) Size() uint32 {
	return uint32(len(f.insts) * 4)
}

// For func links, offset needs to be aligned with RegSize (32 bit=4 bytes)
// uint32 offset is used to record both link position and fill method
// offest/4 is the position of the symbol to be linked in insts[]
// offset%4 is used to record the fill methods

// Constant fill-later methods.
const (
	FillNone = iota // no fill
	FillLink        // fill as linking offset for jump instructions
	FillLow         // fill with the lower 16 bits
	FillHigh        // fill with the higher 16 bits
)

// AddLink links the last instruction in inst to the symbol pkg.sym, where pkg
// and sym are using the indexing of the object file.  fill field must be less
// than 4 so that it fits in the lowest 2 bits in the offset field. The other
// bits of the offset fields will be automatically calculated based on the
// number of instructions in insts.
func (f *Func) AddLink(fill int, ps *PkgSym) {
	if ps.Pkg == "" {
		panic("empty package")
	}

	if len(f.insts) == 0 {
		panic("no inst to link")
	}
	if !(fill > 0 && fill <= 3) {
		panic("invalid fill")
	}

	offset := uint32(len(f.insts))*4 - 4
	offset |= uint32(fill) & 0x3
	link := &link{offset: offset, PkgSym: ps}
	f.links = append(f.links, link)
}
