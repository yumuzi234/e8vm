package arch

type pageOffset struct {
	*page
	offset uint32
}

func (p *pageOffset) readByte(offset uint32) byte {
	return p.page.ReadByte(p.offset + offset)
}

func (p *pageOffset) writeByte(offset uint32, b byte) {
	p.page.WriteByte(p.offset+offset, b)
}

func (p *pageOffset) writeWord(offset uint32, w uint32) {
	p.page.WriteWord(p.offset+offset, w)
}

func (p *pageOffset) readWord(offset uint32) uint32 {
	return p.page.ReadWord(p.offset + offset)
}
