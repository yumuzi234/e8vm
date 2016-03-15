package ir

// Symbol is a linking symbol
type Symbol struct {
	Pkg, Name string
}

// HeapSym is a symbol that lives on the heap
type HeapSym struct {
	*Symbol
	*Attribute
}
