package pl

import (
	"testing"

	"fmt"

	"shanhu.io/smlvm/arch"
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

	oo := func(code, input string) {
		_, es, _ := CompileSingle("main.g", input, false)
		errNum := len(es)
		if errNum != 1 {
			fmt.Println(len(es))
			fmt.Println(input)
		}
		if es == nil {
			t.Log(input)
			t.Error("should error:", code)
			return
		}

		if es[0].Code != code {
			t.Log(input)
			t.Error("ErrCode expected:", code,
				"\nErrcode get:", es[0].Code)
			return
		}
	}

	// Test function for declConflict
	c := func(code, input string) {
		_, es, _ := CompileSingle("main.g", input, false)
		if es == nil || len(es) != 2 {
			t.Log(input)
			t.Error("should have 2 errors for confliction")
			return
		}

		if es[0].Code != code {
			t.Log(input)
			t.Error("ErrCode expected:", code,
				"\nErrcode get:", es[0].Code)
			return
		}

		if es[1].Code != "pl.declConflict.previousPos" {
			t.Log(input)
			t.Error("ErrCode expected:", "pl.declConflict.previousPos",
				"\nErrcode get:", es[1].Code)
			return
		}
	}

	// no main
	oo("pl.missingMainFunc", "")

	// missing returns
	oo("pl.missingFuncReturn", `func f() int { }`)
	oo("pl.missingFuncReturn", `func f() int { for { break } }`)
	oo("pl.missingFuncReturn", `func f() int { for { if true { break } } }`)
	oo("pl.missingFuncReturn", `func f() int { for { if true break } }`)
	oo("pl.missingFuncReturn",
		`func f() int { for true { if true { return 0 } } }`)
	oo("pl.missingFuncReturn",
		`func f() int { for true { if true return 0 } }`)
	oo("pl.missingFuncReturn", `func f() int { if true { return 0 } }`)
	oo("pl.missingFuncReturn", `func f() int { if true return 0 }`)

	// confliction errors return 2 errors,same error code with different pos
	c("pl.declConflict.func",
		`func a() {}; func a() {}`)
	c("pl.declConflict.field",
		`struct A { b int; b int }`)
	c("pl.declConflict.const",
		`const a=1; const a=2`)

	// unused vars
	oo("pl.unusedFuncOrVarible", `func main() { var a int }`)
	oo("pl.unusedFuncOrVarible", `func main() { var a = 3 }`)
	oo("pl.unusedFuncOrVarible", `func main() { a := 3 }`)
	oo("pl.unusedFuncOrVarible", `func main() { var a int; a = 3 }`)
	oo("pl.unusedFuncOrVarible", `func main() { var a int; (a) = (3) }`)
	oo("pl.unusedFuncOrVarible", `func main() { var a, b = 3, 4; _ := a }`)

	// parser, import related
	oo("pl.multiImport", `import(); import()`)

	//expect ';', got keyword
	o("import() func main(){}")

	o(`struct A { a A };`)
	o(`struct A { b B }; struct B { a A };`)
	o(`struct A { b B }; struct B { a [3]A };`)
	o(`struct A { b B }; struct B { a [0]A };`)
	o(`struct A {}; func main() { a := A }`)

	o(`	struct A { func f(){} };
		func main() { var a A; var f func()=a.f; _:=f }`)
	o(`	struct A { func f(){} };
		func main() { var a A; var f func(); f=a.f; _:=f }`)

	o(`struct A { func f(){} }; func main() { var a A; f := a.f; _ := f }`)

	o(` func r() (int, int) { return 3, 4 }
		func p(a, b, c int) { }
		func main() { p(r(), 5) }`)

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
		if !arch.IsErr(e, arch.ErrPanic) {
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
