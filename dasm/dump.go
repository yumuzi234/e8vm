package dasm

import (
	"fmt"
	"io"

	"shanhu.io/smlvm/image"
)

// DumpImage disassembles an image.
func DumpImage(r io.ReadSeeker, out io.Writer) error {
	secs, err := image.Read(r)
	if err != nil {
		return err
	}

	for _, sec := range secs {
		switch sec.Type {
		case image.Code:
			fmt.Fprintln(out, "[code section]")
			lines := Dasm(sec.Bytes, sec.Addr)
			for _, line := range lines {
				fmt.Fprintln(out, line)
			}
		case image.Data:
			fmt.Fprintf(out, "[data of %d bytes at %08x]\n",
				sec.Size, sec.Addr,
			)
			lines := Dasm(sec.Bytes, sec.Addr)
			for _, line := range lines {
				fmt.Fprintln(out, line)
			}
		case image.Zeros:
			fmt.Fprintf(out, "[zeros of %d bytes at %08x]\n",
				sec.Size, sec.Addr,
			)
		case image.Debug:
			fmt.Fprintf(out, "[debug of %d bytes]\n", sec.Size)
		}
	}

	return nil
}
