package g8

import (
	"testing"

	"e8vm.io/e8vm/arch8"
)

func TestSingleFileBad(t *testing.T) {
	o := func(input string) {
		_, es, _ := CompileSingle("main.g", input, false)
		if es == nil {
			t.Log(input)
			t.Error("should error")
			return
		}
	}

	o("") // no main
	o(`func a() {}; func a() {}; func main {}`)
	o(`struct A { b int; b int }; func main {}`)
	o(`struct A { a A }; func main() {}`)
	o(`struct A { b B }; struct B { a A }; func main() {}`)
	o(`struct A { b B }; struct B { a [3]A }; func main() {}`)
	o(`struct A { b B }; struct B { a [0]A }; func main() {}`)
	o(`struct A {}; func main() { a := A }`)

	o(`	struct A { func f(){} }; 
		func main() { var a A; var f func()=a.f; _:=f }`)
	o(`	struct A { func f(){} }; 
		func main() { var a A; var f func(); f=a.f; _:=f }`)

	o(`import(); import()`)
	o("import() func main()")
	o(`struct A { func f(){} }; func main() { var a A; f := a.f; _ := f }`)

	o(` func r() (int, int) { return 3, 4 }
		func p(a, b, c int) { }
		func main() { p(r(), 5) }`)

	// missing returns
	o(`func f() int { }; func main() { }`)
	o(`func f() int { for { break } }; func main() { }`)
	o(`func f() int { for { if true { break } } }; func main() { }`)
	o(`func f() int { for { if true break } }; func main() { }`)
	o(`func f() int { for true { if true { return 0 } } }; func main() { }`)
	o(`func f() int { for true { if true return 0 } }; func main() { }`)
	o(`func f() int { if true { return 0 } }; func main() { }`)
	o(`func f() int { if true return 0 }; func main() { }`)

	// Bugs found by the fuzzer in the past
	o("func main() {}; func f() **o.o {}")
	o("func main() {}; func n()[char[:]]string{}")
	o("func main() {}; func n() { var r = len; _ := r}")
	o("func main() {}; func n() { r := len; _ := r }")
	o("func main() {}; struct A{}; struct A{}")

	o("var a int; func a() {}; func main() {}")
	o("func main() {}; func main() {};")
	o("const a, b = a, b; func main() {}")
	o("const c, d = d, t; func main() {}")
	o(`	func main() {
			var s string
			for i := 0; i < len(s-2); i++ {}
		}`)
}

func TestSingleFilePanic(t *testing.T) {
	// runtime errors

	const N = 100000
	o := func(input string) {
		_, e := singleTestRun(t, input, N)
		if !arch8.IsErr(e, arch8.ErrPanic) {
			t.Log(input)
			t.Log(e)
			t.Error("should panic")
			return
		}
	}

	o("func main() { panic() }")
	o("func main() { var pa *int; printInt(*pa) }")
	o("struct A { a int }; func main() { var pa *A; b := pa.a; _ := b }")
	o("func main() { var a func(); a() }")
	o("func f() {}; func main() { var a func()=f; a=nil; a() }")
	o("func f(p *int) { printInt(*p) }; func main() { f(nil) }")
	o("struct A { p *int }; func main() { var a A; a.p=nil; *a.p=0 }")
}
