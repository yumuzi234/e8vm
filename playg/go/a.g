package main

func p(a []uint) {}

func main() {
	var t [10]uint
	t2 := t[:]
	
	before := uint(&t2[0])
	p(t2[2:5])
	after := uint(&t2[0])
	if before != after {
		panic()
	}
}
