package debug8

import (
	"sort"
)

type funcLocation struct {
	name  string
	start uint32
}

type byStart []*funcLocation

func (f byStart) Len() int { return len(f) }

func (f byStart) Swap(i int, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f byStart) Less(i int, j int) bool {
	return f[i].start < f[j].start
}

func sortTable(t *Table) []*funcLocation {
	var funcs []*funcLocation
	for name, f := range t.Funcs {
		start := f.Start
		var curr funcLocation
		curr.name = name
		curr.start = start
		funcs = append(funcs, &curr)
	}
	sort.Sort(byStart(funcs))
	return funcs
}

func findFunc(pc uint32, funcs []*funcLocation, t *Table) (string, *Func) {
	if len(funcs) == 0 {
		return "", nil
	}

	left := 0
	right := len(funcs) - 1

	for left < right-1 {
		mid := left + (right-left)/2
		if funcs[mid].start == pc {
			return funcs[mid].name, t.Funcs[funcs[mid].name]
		}
		if funcs[mid].start > pc {
			right = mid
		} else {
			left = mid
		}
	}

	if funcs[right].start <= pc {
		f := t.Funcs[funcs[right].name]
		if pc > f.Start+f.Size {
			return "", nil
		}
		return funcs[right].name, f
	}
	f := t.Funcs[funcs[left].name]
	if pc > f.Start+f.Size {
		return "", nil
	}
	return funcs[left].name, f
}
