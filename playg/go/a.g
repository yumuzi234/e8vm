package main

type PageTable struct {
	entries [8]uint
}

func (pt *PageTable) SubTable(i uint) *PageTable {
	return (*PageTable)(pt.Page(i))
}

func (pt *PageTable) Page(i uint) uint {
	return 9
}

func main() {
	var t PageTable
	printUint(t.Page(0))
}
