// Package e8 defines the file format that saves an executable file.
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
