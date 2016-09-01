package debug

import (
	"bytes"
	"fmt"

	"e8vm.io/e8vm/lexing"
)

// Func saves the debug information of a function
type Func struct {
	// compiler filled information
	Frame uint32
	Pos   *lexing.Pos

	// linker filled information
	Start uint32
	Size  uint32
}

func (f *Func) String(name string) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%8x +%4d: ", f.Start, f.Size)
	fmt.Fprintf(buf, "%s", name)
	if f.Pos != nil {
		fmt.Fprintf(buf, "  // %s", f.Pos)
	}
	if f.Frame > 0 {
		fmt.Fprintf(buf, " (frame=%d)", f.Frame)
	}
	return buf.String()
}
