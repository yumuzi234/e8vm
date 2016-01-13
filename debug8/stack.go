package debug8

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
)

func sortTable(t *Table) ([]uint32, []string) {
	var res []string
	var sta []uint32

	for name, f := range t.Funcs {
		start := f.Start
		res = append(res, name)
		sta = append(sta, start)
		i := len(res) - 2
		for i >= 0 && sta[i] > start {
			res[i+1] = res[i]
			res[i] = name
			sta[i+1] = sta[i]
			sta[i] = start
			i--
		}
	}
	return sta, res
}

func findFunc(pc uint32, names []string,
	starts []uint32, t *Table) (string, *Func) {

	left := 0
	right := len(names) - 1

	for left < right-1 {
		mid := left + (right-left)/2
		if starts[mid] == pc {
			return names[mid], t.Funcs[names[mid]]
		}
		if starts[mid] > pc {
			right = mid
		} else {
			left = mid
		}
	}

	if starts[right] <= pc {
		f := t.Funcs[names[right]]
		if pc > f.Start+f.Size {
			return "", nil
		}
		return names[right], f
	}
	f := t.Funcs[names[left]]
	if pc > f.Start+f.Size {
		return "", nil
	}
	return names[left], f
}

// FprintStack prints the stack trace of the current running point.
func FprintStack(w io.Writer, m *arch8.Machine, core byte, t *Table) error {
	regs := m.DumpRegs(core)

	pc := regs[arch8.PC]
	sp := regs[arch8.SP]
	ret := regs[arch8.RET]

	level := 0
	starts, names := sortTable(t)

	for {
		level++

		name, f := findFunc(pc, names, starts, t)
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

		if f.Size <= 4 || f.Frame == 0 { // cannot be a normal function
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
