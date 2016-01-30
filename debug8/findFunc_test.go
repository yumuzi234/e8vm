package debug8

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestFindFunc(t *testing.T) {
	as := func(cond bool, s string, args ...interface{}) {
		if !cond {
			t.Fatalf(s, args...)
		}
	}

	eo := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}

	tbl := NewTable()
	var sum uint32

	for i := 0; i < 10; i++ {
		name := strconv.Itoa(i)
		size := rand.Uint32()
		for size > 10000 {
			size = rand.Uint32()
		}
		tbl.Funcs[name] =
			&Func{Size: size, Start: sum}

		sum = sum + size + 1
	}

	pc := rand.Uint32()
	for pc > sum {
		pc = rand.Uint32()
	}

	fLoc := sortTable(tbl)
	_, f := findFunc(pc, fLoc, tbl)

	as(f != nil, "did not find file")
	eo(f.Start >= pc || f.Size+f.Start <= pc,
		"pc is %v\n start is %v and length is %v\n", pc, f.Start, f.Size)

}
