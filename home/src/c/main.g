import (
	"b"
	"fmt"
)

func main() {
	printInt(555)
	fmt.PrintUint(333)

	f()
}

func f() { f1() }
func f1() { f2() }

func f2() {
	var a *int
	var b = *a
	panic()
}
