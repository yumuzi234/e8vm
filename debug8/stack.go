package debug8

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
)

func findFunc(pc uint32, t *Table) (string, *Func) {
	for name, f := range t.Funcs {
		if pc >= f.Start && pc < f.Start+f.Size {
			return name, f
		}
	}
	return "", nil
}

// FprintStack prints the stack trace of the current running point.
func FprintStack(w io.Writer, m *arch8.Machine, core byte, t *Table) error {
	regs := m.DumpRegs(core)

	pc := regs[arch8.PC]
	sp := regs[arch8.SP]
	ret := regs[arch8.RET]

	level := 0

	for {
		level++

		name, f := findFunc(pc, t)
		if f == nil {
			if level == 1 {
				_, err := fmt.Fprintf(w, "? pc=%08x\n", pc)
				return err
			}
			return nil
		}

		_, err := fmt.Fprintln(w, funcString(name, f))
		if err != nil {
			return err
		}

		if f.Size <= 4 { // cannot be a normal function
			if level != 1 {
				// calling in a non-normal function
				// return
				return nil
			}

			// use ret as pc
			pc = ret
			continue
		}

		retAddr, err := m.ReadWord(core, sp+f.Frame-4)
		if err != nil {
			_, err := fmt.Fprintf(w, "! unable to recover: %s\n", err)
			return err
		}

		pc = retAddr
		sp = sp + f.Frame
	}
}
