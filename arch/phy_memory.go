package arch

import (
	"io"
	"math"
)

// PhyMemory is a collection of contiguous pages.
type phyMemory struct {
	npage uint32
	pages map[uint32]*page
}

// NewPhyMemory creates a physical memory of size bytes.
func newPhyMemory(size uint32) *phyMemory {
	if size%PageSize != 0 {
		panic("size misaligned")
	}

	ret := new(phyMemory)
	ret.pages = make(map[uint32]*page)
	if size > 0 {
		ret.npage = size / PageSize
	} else {
		ret.npage = (math.MaxUint32 + 1) / PageSize
	}
	if ret.npage < 16 {
		ret.npage = 16
	}

	return ret
}

// Size returns the size of the physical memory.
func (pm *phyMemory) Size() uint32 {
	return pm.npage * PageSize
}

// Page returns the page for the particular page number
// Returns nil when the page number is out of range
func (pm *phyMemory) Page(pn uint32) *page {
	if pn == 0 || pn >= pm.npage {
		return nil // out of range
	}

	ret, found := pm.pages[pn]
	if !found {
		// create an empty page on demand
		ret = newPage()
		pm.pages[pn] = ret
	}

	return ret
}

func (pm *phyMemory) pageForU8(addr uint32) (*page, *Excep) {
	p := pm.Page(addr / PageSize)
	if p == nil {
		return nil, newOutOfRange(addr)
	}
	return p, nil
}

func (pm *phyMemory) pageForU32(addr uint32) (*page, *Excep) {
	if addr%4 != 0 {
		return nil, errMisalign
	}
	return pm.pageForU8(addr)
}

// ReadU8 reads the byte at the given address.
// If the address is out of range, it returns an error.
func (pm *phyMemory) ReadU8(addr uint32) (byte, *Excep) {
	p, e := pm.pageForU8(addr)
	if e != nil {
		return 0, e
	}
	return p.ReadU8(addr), nil
}

// WriteU8 writes the byte at the given address.
// If the address is out of range, it returns an error.
func (pm *phyMemory) WriteU8(addr uint32, v byte) *Excep {
	p, e := pm.pageForU8(addr)
	if e != nil {
		return e
	}
	p.WriteU8(addr, v)
	return e
}

// ReadU32 reads the byte at the given address.
// If the address is out of range or not 4-byte aligned, it returns an error.
func (pm *phyMemory) ReadU32(addr uint32) (uint32, *Excep) {
	p, e := pm.pageForU32(addr)
	if e != nil {
		return 0, e
	}
	return p.ReadU32(addr), nil
}

// WriteU32 reads the byte at the given address.
// If the address is out of range or not 4-byte aligned, it returns an error.
func (pm *phyMemory) WriteU32(addr uint32, v uint32) *Excep {
	p, e := pm.pageForU32(addr)
	if e != nil {
		return e
	}
	p.WriteU32(addr, v)
	return nil
}

func (pm *phyMemory) writeBytes(r io.Reader, offset uint32) error {
	start := offset % PageSize
	pageBuf := make([]byte, PageSize)
	pn := offset / PageSize
	for {
		p := pm.Page(pn)
		if p == nil {
			return newOutOfRange(offset)
		}

		buf := pageBuf[:PageSize-start]
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}

		p.WriteAt(buf[:n], start)
		start = 0
		pn++
	}

	return nil
}
