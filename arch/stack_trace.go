package arch

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/image"
)

type funcEntry struct {
	name  string
	start uint32
}

type byStart []*funcEntry

func (f byStart) Len() int { return len(f) }

func (f byStart) Swap(i int, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f byStart) Less(i int, j int) bool {
	return f[i].start < f[j].start
}

func sortTable(t *debug.Table) []*funcEntry {
	var funcs []*funcEntry
	for name, f := range t.Funcs {
		funcs = append(funcs, &funcEntry{
			name:  name,
			start: f.Start,
		})
	}
	sort.Sort(byStart(funcs))
	return funcs
}

func findFunc(fs []*funcEntry, pc uint32, t *debug.Table) (
	string, *debug.Func,
) {
	if len(fs) == 0 {
		return "", nil
	}

	left := 0
	right := len(fs) - 1

	for left < right-1 {
		mid := left + (right-left)/2
		if fs[mid].start == pc {
			return fs[mid].name, t.Funcs[fs[mid].name]
		}
		if fs[mid].start > pc {
			right = mid
		} else {
			left = mid
		}
	}

	if fs[right].start <= pc {
		f := t.Funcs[fs[right].name]
		if pc > f.Start+f.Size {
			return "", nil
		}
		return fs[right].name, f
	}
	f := t.Funcs[fs[left].name]
	if pc <= f.Start || pc > f.Start+f.Size {
		return "", nil
	}
	return fs[left].name, f
}

func debugSection(secs []*image.Section) *image.Section {
	for _, sec := range secs {
		if sec.Type == image.Debug {
			return sec
		}
	}
	return nil
}

// FprintStack prints the stack trace of a machine from its exception
// and registers.
func FprintStack(w io.Writer, m *Machine, excep *CoreExcep) error {
	sec := debugSection(m.Sections)
	if sec == nil {
		return errors.New("debug section not found")
	}

	t, err := debug.UnmarshalTable(sec.Bytes)
	if err != nil {
		return err
	}
	funcs := sortTable(t)

	core := byte(excep.Core)
	regs := m.DumpRegs(core)
	pc := regs[PC]
	sp := regs[SP]
	ret := regs[RET]
	level := 0

	fmt.Fprintf(w, "err: %s\n", excep.Err.Error())
	fmt.Fprintf(w, "core=%d excep=%d\n", core, excep.Code)
	fmt.Fprintf(w, "pc=%08x sp=%08x ret=%08x\n", pc, sp, ret)
	inst, readErr := m.ReadWord(0, pc)
	if readErr == nil {
		fmt.Fprintf(w, "inst=%08x\n", inst)
	}

	for {
		level++

		name, f := findFunc(funcs, pc, t)
		if f == nil {
			if level == 1 {
				_, err := fmt.Fprintf(w, "? pc=%08x\n", pc)
				return err
			}
			return nil
		}

		_, err := fmt.Fprintln(w, f.String(name))
		if err != nil {
			return err
		}

		if f.Size <= 4 || f.Frame == 0 { // cannot be a normal function
			if level != 1 {
				// calling in a non-normal function
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
