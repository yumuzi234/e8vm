package debug8

import (
	"sort"
)

type funcLocation struct {
	name  string
	start uint32
}

type funcLocationSlice []*funcLocation

func (f funcLocationSlice) Len() int {
	return len(f)
}

func (f funcLocationSlice) Swap(i int, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f funcLocationSlice) Less(i int, j int) bool {
	return f[i].start < f[j].start
}

func sortTable(t *Table) funcLocationSlice {

	var fLocation funcLocationSlice
	for name, f := range t.Funcs {
		start := f.Start
		var curr funcLocation
		curr.name = name
		curr.start = start
		fLocation = append(fLocation, &curr)
	}
	sort.Sort(fLocation)
	return fLocation
}

func findFunc(pc uint32, fLoc funcLocationSlice, t *Table) (string, *Func) {
	if len(fLoc) == 0 {
		return "", nil
	}

	left := 0
	right := len(fLoc) - 1

	for left < right-1 {
		mid := left + (right-left)/2
		if fLoc[mid].start == pc {
			return fLoc[mid].name, t.Funcs[fLoc[mid].name]
		}
		if fLoc[mid].start > pc {
			right = mid
		} else {
			left = mid
		}
	}

	if fLoc[right].start <= pc {
		f := t.Funcs[fLoc[right].name]
		if pc > f.Start+f.Size {
			return "", nil
		}
		return fLoc[right].name, f
	}
	f := t.Funcs[fLoc[left].name]
	if pc > f.Start+f.Size {
		return "", nil
	}
	return fLoc[left].name, f
}
