package arch8

import (
	"math/rand"
	"strconv"
	"testing"

	"e8vm.io/e8vm/debug8"
)

func TestFindFunc(t *testing.T) {
	eo := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}

	tbl := debug8.NewTable()
	funcs := sortTable(tbl)
	_, f := findFunc(funcs, 1, tbl)
	eo(f != nil, "func not found")

	tbl = &debug8.Table{
		Funcs: map[string]*debug8.Func{
			"f1": {Size: 20, Start: 0},
		},
	}
	funcs = sortTable(tbl)
	_, f = findFunc(funcs, 30, tbl)
	eo(f != nil, "findFunc error")
	_, f = findFunc(funcs, 10, tbl)
	eo(f.Start != 0, "findFunc error")

	tbl = &debug8.Table{
		Funcs: map[string]*debug8.Func{
			"f1": {Size: 20, Start: 0},
			"f2": {Size: 1, Start: 25},
			"f3": {Size: 20, Start: 27},
		},
	}
	funcs = sortTable(tbl)
	_, f = findFunc(funcs, 26, tbl)
	eo(f.Start != 25,
		"pc is %v\n start is %v and length is %v\n", 25, f.Start, f.Size)
	_, f = findFunc(funcs, 24, tbl)
	eo(f != nil, "findFunc error")
	_, f = findFunc(funcs, 50, tbl)
	eo(f != nil, "findFunc error")

	tbl = debug8.NewTable()
	var sum uint32
	for i := 0; i < 100; i++ {
		name := "f" + strconv.Itoa(i)
		size := rand.Uint32()
		tbl.Funcs[name] = &debug8.Func{Size: size, Start: sum}
		sum = sum + size + 1
	}
	pc := rand.Uint32()
	for pc > sum {
		pc = rand.Uint32()
	}
	funcs = sortTable(tbl)
	_, f = findFunc(funcs, pc, tbl)

	eo(f.Start >= pc || f.Size+f.Start <= pc,
		"pc is %v\n func starts %v and length is %v\n", pc, f.Start, f.Size)
}
