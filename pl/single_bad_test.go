package pl

import (
	"testing"

	"shanhu.io/smlvm/arch"
)

func TestSingleFileBad(t *testing.T) {
	oo := func(code, input string) {
		_, es, _ := CompileSingle("main.g", input, false)
		errNum := len(es)
		if errNum != 1 {
			t.Log(len(es))
			for i := 0; i < len(es); i++ {
				t.Log(es[i].Code)
			}
			t.Log(input)
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

	// confliction return 2 errors
	c("declConflict.func", `func a() {}; func a() {}`)
	c("declConflict.field", `struct A { b int; b int }`)
	c("declConflict.const", `const a=1; const a=2`)
	c("declConflict.struct", "struct A{}; struct A{}")
	c("declConflict.Var", "var a int; func a() {}")
	c("declConflict.func", "func main() {}; func main() {};")

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
	oo("missingSemi", "import() func main(){}")

	// circular dependence
	oo("circDep.struct", `struct A { a A };`)
	oo("circDep.struct", `struct A { b B }; struct B { a A };`)
	oo("circDep.struct", `struct A { b B }; struct B { a [3]A };`)
	oo("circDep.struct", `struct A { b B }; struct B { a [0]A };`)
	oo("circDep.const", `const a = b; const b = a`)
	oo("circDep.const", `const a = 3 + b; const b = a`)
	oo("circDep.const", `const a = 3 + b; const b = 0 - a`)
	oo("circDep.const", "const a, b = a, b")

	//Assign and allocate
	oo("cannotAlloc", `struct A {}; func main() { a := A }`)
	oo("cannotAlloc", `struct A {}; func (a *A) f(){};
		func main() { var a A; f := a.f; _ := f }`)
	oo("cannotAlloc", "func n() { var r = len; _ := r}")
	oo("cannotAlloc", "func n() { r := len; _ := r }")

	oo("cannotAssign", `struct A {}; func (a *A) f(){};
		func main() { var a A; var f func()=a.f; _:=f }`)
	oo("cannotAssign", `struct A {}; func (a *A) f(){};
		func main() { var a A; var f func(); f=a.f; _:=f }`)
	oo("cannotAssign", `func main() { var a [2]int; var b [3]int;
		a=b}`)
	// If the length are not the same, cannot assign either and it will be
	// another error code

	oo("multiRefInExprList", ` func r() (int, int) { return 3, 4 }
		func p(a, b, c int) { }
		func main() { p(r(), 5) }`)

	oo("elseStart", `func main() {
		if true { }
		else { } }`)

	// Bugs found by the fuzzer in the past
	oo("undefinedIdent", "func f() **o.o {}")
	oo("expectConstExpr", "func n()[char[:]]string{}")
	oo("undefinedIdent", "const c, d = d, t; func main() {}")

	oo("cannotCast", `func main() {
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
