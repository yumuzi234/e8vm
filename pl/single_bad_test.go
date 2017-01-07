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
			for i := 0; i < len(es); i++ {
				fmt.Println(es[i].Code)
			}
			fmt.Println(input)
		}
		if es == nil {
			t.Log(input)
			t.Error("should error:", code)
			return
		}
		code = "pl." + code
		if es[0].Code != code {
			t.Log(input)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
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
		code = "pl." + code
		if es[0].Code != code {
			t.Log(input)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
			return
		}

		const previousPos = "pl.declConflict.previousPos"
		if es[1].Code != previousPos {
			t.Log(input)
			t.Error("ErrCode expected: %q, got %q", previousPos, es[1].Code)
			return
		}
	}

	// no main
	oo("missingMainFunc", "")

	// missing returns
	oo("missingReturn", `func f() int { }`)
	oo("missingReturn", `func f() int { for { break } }`)
	oo("missingReturn", `func f() int { for { if true { break } } }`)
	oo("missingReturn", `func f() int { for { if true break } }`)
	oo("missingReturn",
		`func f() int { for true { if true { return 0 } } }`)
	oo("missingReturn",
		`func f() int { for true { if true return 0 } }`)
	oo("missingReturn", `func f() int { if true { return 0 } }`)
	oo("missingReturn", `func f() int { if true return 0 }`)

	// confliction errors return 2 errors, same error code with different pos
	c("declConflict.func", `func a() {}; func a() {}`)
	c("declConflict.field", `struct A { b int; b int }`)
	c("declConflict.const", `const a=1; const a=2`)
	c("declConflict.struct",
		"func main() {}; struct A{}; struct A{}")

	// unused vars
	oo("unusedSym", `func main() { var a int }`)
	oo("unusedSym", `func main() { var a = 3 }`)
	oo("unusedSym", `func main() { a := 3 }`)
	oo("unusedSym", `func main() { var a int; a = 3 }`)
	oo("unusedSym", `func main() { var a int; (a) = (3) }`)
	oo("unusedSym", `func main() { var a, b = 3, 4; _ := a }`)

	// parser, import related
	oo("multiImport", `import(); import()`)

	// expect ';', got keyword
	oo("missingSemicolon", "import() func main(){}")

	//circular dependence
	oo("circDep.struct", `struct A { a A };`)
	oo("circDep.struct", `struct A { b B }; struct B { a A };`)
	oo("circDep.struct", `struct A { b B }; struct B { a [3]A };`)
	oo("circDep.struct", `struct A { b B }; struct B { a [0]A };`)
	oo("circDep.const", `const a = b; const b = a`)
	oo("circDep.const", `const a = 3 + b; const b = a`)
	oo("circDep.const", `const a = 3 + b; const b = 0 - a`)

	oo("CannotAllocte", `struct A {}; func main() { a := A }`)
	oo("CannotAllocte", `struct A {}; func (a *A) f(){};
		func main() { var a A; f := a.f; _ := f }`)
	oo("CannotAllocte",
		"func main() {}; func n() { var r = len; _ := r}")
	oo("CannotAllocte", "func main() {}; func n() { r := len; _ := r }")

	//If the len is the same, cannot assign either
	//and it will be another error code
	oo("CannotAssign", `struct A {}; func (a *A) f(){};
		func main() { var a A; var f func()=a.f; _:=f }`)
	oo("CannotAssign", `struct A {}; func (a *A) f(){};
		func main() { var a A; var f func(); f=a.f; _:=f }`)

	oo("multiRefInExprList", ` func r() (int, int) { return 3, 4 }
		func p(a, b, c int) { }
		func main() { p(r(), 5) }`)

	oo("undefinedIdent", "func main() {}; func f() **o.o {}")
	oo("expConstExpr", "func main() {}; func n()[char[:]]string{}")

	// Bugs found by the fuzzer in the past
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
