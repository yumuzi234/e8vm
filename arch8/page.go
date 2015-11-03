package arch8

// PageSize is the number of bytes a page contains.
const PageSize = 4096

// Page is a memory addressable area of PageSize bytes
type page struct {
	bs []byte
}

// NewPage creates a new empty page.
func newPage() *page {
	return &page{
		bs: make([]byte, PageSize),
	}
}

// ReadByte reads a byte at the particular offset.
// When offset is larger than offset, it uses the modular.
func (p *page) ReadByte(offset uint32) byte {
	return p.bs[offset%PageSize]
}

// WriteByte writes a byte into the page at a particular offset.
// When offset is larger than offset, it uses the modular.
func (p *page) WriteByte(offset uint32, b byte) {
	p.bs[offset%PageSize] = b
}

func wordOff(offset uint32) uint32 {
	return (offset % PageSize) & ^uint32(0x3)
}

// ReadWord reads the word at the particular offset.
// When offset is larger than offset, it uses the modular.
// When offset is not 4-byte aligned, it aligns down.
func (p *page) ReadWord(offset uint32) uint32 {
	offset = wordOff(offset)
	return Endian.Uint32(p.bs[offset : offset+4])
}

// WriteWord writes the word at the particular offset.
// When offset is larger than offset, it uses the modular.
// When offset is not 4-byte aligned, it aligns down.
func (p *page) WriteWord(offset uint32, w uint32) {
	offset = wordOff(offset)
	Endian.PutUint32(p.bs[offset:offset+4], w)
}

// WriteAt writes a series of bytes starting at offset
func (p *page) WriteAt(bs []byte, offset uint32) {
	for i, b := range bs {
		off := offset + uint32(i)
		if off > PageSize {
			panic("out of range")
		}
		p.WriteByte(off, b)
	}
}
