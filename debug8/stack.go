package debug8

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
)

// FprintStack prints the stack trace of the current running point.
func FprintStack(w io.Writer, m *arch8.Machine, core byte, t *Table) error {
	regs := m.DumpRegs(core)

	pc := regs[arch8.PC]
	sp := regs[arch8.SP]
	ret := regs[arch8.RET]

	_, _ = sp, ret

	for name, f := range t.Funcs {
		if pc >= f.Start && pc < f.Start+f.Size {
			_, err := fmt.Fprintln(w, funcString(name, f))
			return err
		}
	}

	fmt.Fprintln(w, "? pc=%08x", pc)
	return nil
}
