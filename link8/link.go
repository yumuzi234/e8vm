package link8

// link is symbol to be linked in an instruction in a code section at a
// particular offset, or in a piece of data in a data section at a particular
// offset.  it uses the indices in the package for symbol lookup
type link struct {
	offset uint32
	pkg    uint32 // relative package index
	sym    uint32
}
