package e8

// Section types
const (
	None uint16 = iota
	Code
	Data
	Zeros // a.k.a. BSS
	Symbols
	DebugInfo
	Comment
)
