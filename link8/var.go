package link8

import (
	"bytes"
	"fmt"
	"math"

	"e8vm.io/e8vm/arch8"
)

// Var reprsents a variable object
type Var struct {
	align uint32
	buf   *bytes.Buffer

	addr   uint32
	prePad uint32

	zeros uint32

	links []*link // symbols
}

// NewVar creates a new relocatable data section.
func NewVar(align uint32) *Var {
	ret := new(Var)

	if align == 0 {
		align = 1
	} else if align != 1 && align != 4 {
		panic("invalid align")
	}

	ret.align = align
	ret.buf = new(bytes.Buffer)

	return ret
}

// Zeros set this variable as a BSS section.
func (v *Var) Zeros(n uint32) { v.zeros = n }

// IsZeros checks if the variable section is a BSS section.
func (v *Var) IsZeros() bool { return v.zeros > 0 }

// Write appends bytes to this data section.
func (v *Var) Write(buf []byte) (int, error) {
	if v.zeros != 0 {
		panic("cannot write to zeros section")
	}
	return v.buf.Write(buf)
}

// Pad pads n bytes into this data section
func (v *Var) Pad(n uint32) {
	if v.zeros != 0 {
		panic("cannot pad to zeros section")
	}
	v.buf.Write(make([]byte, n))
}

// WriteLink writes a symbol link into the data section.
func (v *Var) WriteLink(pkg, sym string) error {
	if pkg == "" {
		panic("empty package")
	}

	if v.zeros != 0 {
		panic("cannot write link to zeros section")
	}

	if v.align%arch8.RegSize != 0 {
		return fmt.Errorf("align %d, not register size aligned", v.align)
	}
	offset := uint32(v.buf.Len())
	if offset%arch8.RegSize != 0 {
		return fmt.Errorf("offset %d, not register size aligned", offset)
	}

	lnk := &link{
		offset: uint32(v.buf.Len()),
		pkg:    pkg,
		sym:    sym,
	}
	v.links = append(v.links, lnk)
	v.Pad(arch8.RegSize) // symbol has a size of a register
	return nil
}

// Size returns the current size of the section
func (v *Var) Size() uint32 {
	if v.zeros != 0 {
		return v.zeros
	}
	return uint32(v.buf.Len())
}

// TooLarge checks if the size is larger than 2GB
func (v *Var) TooLarge() bool {
	return v.Size() > math.MaxInt32
}
