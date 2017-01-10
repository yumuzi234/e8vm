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
	// I will introduce the function similar as single_bad_test.go
	oo := func(code, s string) {
		buf := strings.NewReader(s)
		_, _, es := File("test.g", buf, false)
		errNum := len(es)
		if errNum != 1 {
			t.Log(len(es))
			for i := 0; i < len(es); i++ {
				t.Log(es[i].Code)
			}
			t.Log(s)
		}
		if es == nil {
			t.Log(s)
			t.Log("should fail")
			t.Fail()
		}
		code = "pl." + code
		if es[0].Code != code {
			t.Log(s)
			t.Log(es[0].Err)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
			return
		}
	}
	var testCases = []struct {
		s    string
		code string
	}{
		{"func f() ", "expectOp"},
		{"func f{} ", "expectOp"},
		{"func f(", "expectOp"},
		{"func f)", "expectOp"},
		{"func f; {}", "expectOp"},
		{"func f(a int) () {}", "expectReturnList"},
		{"func f(,a) {}", "expectType"},
		{"func f(a int) (,a) {}", "expectType"},
		{"func f(a b int) (,a) {}", "expectOp"},
		{"func f(a int,,) (,a) {}", "expectType"},
	}
	// The test part can be fulfilled like that or the same as single_bad_test
	for _, c := range testCases {
		oo(c.code, c.s)
	}

	for _, s := range []string{
		"func f(a int \n b int) {}",
		"func f() \n {}",
		"var (a int, b int); func main() { }",
		"var (a int, b int}",
		"import (); import()",
		`var (a "a"}`,
		`var (a "a";}`,
	} {
		buf := strings.NewReader(s)
		_, _, es := File("test.g", buf, false)
		if es == nil {
			t.Log(s)
			t.Log("should fail")
			t.Fail()
		}
	}
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
