package debug8

import "fmt"
import "testing"
import "strconv"
import "math/rand"

func BenchmarkTest(b *testing.B) {

	for n := 0; n < b.N; n++ {
		test()
	}
}

func test() (uint32, *Func) {

	t := NewTable()
	var starts uint32

	for i := 0; i < 10; i++ {
		name := strconv.Itoa(i)
		frameSize := rand.Uint32()
		for frameSize > 10000 {
			frameSize = rand.Uint32()
		}	

		t.Funcs[name] = 
		&Func{Size: frameSize, Start: starts }

		starts = starts + frameSize + 1
	}

	// for name, f := range t.Funcs {
	// 	fmt.Println(name)
	// 	fmt.Println(f.Start)
	// }

	pc := rand.Uint32()
	for pc > starts {
		pc = rand.Uint32()
	}

	Starts, Names := sortTable(t)
	_, f := findFunc(pc, Names, Starts, t)

	return pc, f

}

func TestFindFunc(t *testing.T) {

	pc, res := test()
	fmt.Printf( "pc is %v\n start is %v and length is %v\n", pc, res.Start, res.Size)

}
