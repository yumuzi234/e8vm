package arch

type pageOffset struct {
	*page
	offset uint32
}

func (p *pageOffset) readU8(offset uint32) byte {
	return p.page.ReadU8(p.offset + offset)
}

func (p *pageOffset) writeU8(offset uint32, b byte) {
	p.page.WriteU8(p.offset+offset, b)
}

func (p *pageOffset) writeU32(offset uint32, w uint32) {
	p.page.WriteU32(p.offset+offset, w)
}

func (p *pageOffset) readU32(offset uint32) uint32 {
	return p.page.ReadU32(p.offset + offset)
}
