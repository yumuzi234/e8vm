package debug8

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func test() (uint32, *Func) {
	t := NewTable()
	var sum uint32

	for i := 0; i < 10; i++ {
		name := strconv.Itoa(i)
		size := rand.Uint32()
		for size > 10000 {
			size = rand.Uint32()
		}
		t.Funcs[name] =
			&Func{Size: size, Start: sum}

		sum = sum + size + 1
	}

	pc := rand.Uint32()
	for pc > sum {
		pc = rand.Uint32()
	}

	starts, names := sortTable(t)
	_, f := findFunc(pc, names, starts, t)

	return pc, f

}

func TestFindFunc(t *testing.T) {
	pc, res := test()
	fmt.Printf("pc is %v\n start is %v and length is %v\n", pc, res.Start, res.Size)

}
