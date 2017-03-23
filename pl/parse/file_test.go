package parse

import (
	"testing"

	"strings"
)

func TestFile_good(t *testing.T) {
	for _, s := range []string{
		"func f() {}",
		"func f() {\n}",
		"func f(int) {}",
		"func f(a int) {}",
		"func f(a int,) {}",
		"func f(a, b int) {}",
		"func f(a int, b int) {}",
		"func f() (int) {}",
		"func f() (a int) {}",
		"func f() (a, b int) {}",
		"func f() (a, b int,) {}",
		"func f() (a int, b int) {}",
		"func f(int) (a int, b int) {}",
		"func f(int) (a int, b int,) {}",
		`func f(int) (
			a int,
			b int,
		) {}
		`,
		`import ( a "a" )`,
		`import ( "a" )`,
		`import ( _ "a" )`,
		`import ( "a"; "b" )`,
		`import ( a "a"; "b" )`,
		`interface T {
			add() int
			print() (a string, b char)
		}`,
		`interface T {
			add(int, int) int
			print(string) 
		}`,
		`interface T {
			divide(a, b int) (int, int)
		}`,
		`interface T {
		}`,
	} {
		buf := strings.NewReader(s)
		f, _, es := File("test.g", buf, false)
		if es != nil {
			t.Log(s)
			for _, e := range es {
				t.Log(e)
			}
			t.Fail()
		} else if f == nil {
			t.Log(s)
			t.Log("returned nil")
			t.Fail()
		}
	}
}

func TestFile_bad(t *testing.T) {
	o := func(code, s string) {
		buf := strings.NewReader(s)
		_, _, es := File("test.g", buf, false)

		if len(es) != 1 {
			t.Log(s)
			t.Logf("%d errors returned", len(es))
			for _, e := range es {
				t.Log(e.Code)
				t.Log(e.Err)
			}
		}

		if es == nil {
			t.Log(s)
			t.Log("should fail")
			t.Fail()
			return
		}

		code = "pl." + code
		if es[0].Code != code {
			t.Log(s)
			t.Log(es[0].Err)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
		}
	}

	o("expectOp", "func f()")
	o("expectOp", "func f{}")
	o("expectOp", "func f(")
	o("expectOp", "func f)")
	o("expectOp", "func f; {}")
	o("expectReturnList", "func f(a int) () {}")
	o("expectType", "func f(,a) {}")
	o("expectType", "func f(a int) (,a) {}")
	o("expectOp", "func f(a b int) (,a) {}")
	o("expectType", "func f(a int,,) (,a) {}")
	o("expectOp", "func f(a int \n b int) {}")
	o("expectOp", "func f() \n {}")
	o("missingSemi", "var (a int, b int)")
	o("multiImport", "import (); import()")
	o("expectType", `var (a "a")`)
	o("expectType", `var (a "a";)`)
	o("unexpected", `struct s{
		func f() {}	
	}`)
}

func TestFileTokens(t *testing.T) {
	buf := strings.NewReader("func f() {}")
	_, rec, es := File("test.g", buf, false)
	if es != nil {
		for _, e := range es {
			t.Log(e)
		}
		t.Fail()
	}
	toks := rec.Tokens()
	for _, tok := range toks {
		t.Log(tok.Pos, tok.Lit)
	}
	firstTok := toks[0]
	pos := firstTok.Pos
	if pos.Line != 1 || pos.Col != 1 {
		t.Error("first token not starting with test.g:1:1")
	}
}
