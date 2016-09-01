// Package image defines the file format that saves an executable file.
package image

// Section types
const (
	None uint8 = iota
	Code
	Data
	Zeros // a.k.a. BSS
	Debug
	Comment
)
