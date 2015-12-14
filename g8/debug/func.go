package e8

// Func saves the debug information of a function
type Func struct {
	Pkg  string
	Name string

	Start uint32
	Size  uint32
	Frame uint32

	File string
	Line int
	Col  int
}
