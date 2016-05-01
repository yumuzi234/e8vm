import (
	"b"
	"fmt"
	
	"asm/test2"
)

func main() {
	printInt(555)
	fmt.PrintUint(333)

	f()
}

func f() { f1() }
func f1() { f2() }

struct A {
	// func f3() = test2.A
}

func f2() {
	var a *int
	var b = *a
	_ := b
	panic()
}
