package main

type PageTable struct {
	entries []uint
}

type a struct {
	_ int
}

func (pt *PageTable) p() {
	var aa a
	aa.p(pt.entries)
}

func (*a) p(a []uint) {
	printUint(uint(&a[0]))
}

func main() {
	var t [10]uint
	var pt PageTable
	var aa a
	pt.entries = t[:]
	printUint(uint(&pt.entries[0]))
	aa.p(pt.entries)
	pt.p()
}
