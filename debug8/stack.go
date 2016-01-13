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

func findFunc(pc uint32, Names []string, Starts []uint32, t *Table) (string, *Func) {
	
	left := 0
	right := len(Names) - 1

	for left < right-1 {
		mid := left + (right-left)/2
		if Starts[mid] == pc {
			return Names[mid], t.Funcs[Names[mid]]
		}
		if Starts[mid] > pc {
			right = mid
		} else {
			left = mid
		}
	}

	if Starts[right] <= pc {
		f := t.Funcs[Names[right]]
		if pc > f.Start + f.Size {
		return "", nil
		}
		return Names[right], f
	} else {
		f := t.Funcs[Names[left]]
		if pc > f.Start + f.Size {
		return "", nil
		}
		return Names[left], f
	}
}

//FprintStack prints the stack trace of the current running point.
func FprintStack(w io.Writer, m *arch8.Machine, core byte, t *Table) error {
	regs := m.DumpRegs(core)

	pc := regs[arch8.PC]
	sp := regs[arch8.SP]
	ret := regs[arch8.RET]

	level := 0
	Starts, Names := sortTable(t)

	for {
		level++

		name, f := findFunc(pc, Names, Starts, t)
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
