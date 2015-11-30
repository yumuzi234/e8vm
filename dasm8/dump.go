package dasm8

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/e8"
)

// DumpImage disassembles an image.
func DumpImage(r io.ReadSeeker, out io.Writer) error {
	secs, err := e8.Read(r)
	if err != nil {
		return err
	}

	for _, sec := range secs {
		switch sec.Type {
		case e8.Code:
			fmt.Fprintln(out, "[code section]")
			lines := Dasm(sec.Bytes, sec.Addr)
			for _, line := range lines {
				fmt.Fprintln(out, line)
			}
		case e8.Data:
			fmt.Fprintf(out, "[data of %d bytes at %08x]\n",
				sec.Size, sec.Addr,
			)
			lines := Dasm(sec.Bytes, sec.Addr)
			for _, line := range lines {
				fmt.Fprintln(out, line)
			}
		case e8.Zeros:
			fmt.Fprintf(out, "[zeros of %d bytes at %08x]\n",
				sec.Size, sec.Addr,
			)
		}
	}

	return nil
}
