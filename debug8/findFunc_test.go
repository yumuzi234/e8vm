package debug8

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestFindFunc(t *testing.T) {
	eo := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}

	tbl := NewTable()
	fLoc := sortTable(tbl)
	_, f := findFunc(1, fLoc, tbl)
	eo(f != nil, "findFunc error")

	tbl = NewTable()
	tbl.Funcs["1"] =
		&Func{Size: 20, Start: 0}
	fLoc = sortTable(tbl)
	_, f = findFunc(30, fLoc, tbl)
	eo(f != nil, "findFunc error")
	_, f = findFunc(10, fLoc, tbl)
	eo(f.Start != 0, "findFunc error")

	tbl = NewTable()
	tbl.Funcs["1"] =
		&Func{Size: 20, Start: 0}
	tbl.Funcs["2"] =
		&Func{Size: 1, Start: 25}
	tbl.Funcs["3"] =
		&Func{Size: 20, Start: 27}
	fLoc = sortTable(tbl)
	_, f = findFunc(26, fLoc, tbl)
	eo(f.Start != 25,
		"pc is %v\n start is %v and length is %v\n", 25, f.Start, f.Size)
	_, f = findFunc(24, fLoc, tbl)
	eo(f != nil, "findFunc error")
	_, f = findFunc(50, fLoc, tbl)
	eo(f != nil, "findFunc error")

	tbl = NewTable()
	var sum uint32
	for i := 0; i < 100; i++ {
		name := strconv.Itoa(i)
		size := rand.Uint32()
		tbl.Funcs[name] =
			&Func{Size: size, Start: sum}

		sum = sum + size + 1
	}
	pc := rand.Uint32()
	for pc > sum {
		pc = rand.Uint32()
	}
	fLoc = sortTable(tbl)
	_, f = findFunc(pc, fLoc, tbl)

	eo(f.Start >= pc || f.Size+f.Start <= pc,
		"pc is %v\n func starts %v and length is %v\n", pc, f.Start, f.Size)
}
