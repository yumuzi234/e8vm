package parse

import (
	"testing"

	"strings"

	"shanhu.io/smlvm/pl/ast"
)

func TestStmts_good(t *testing.T) {
	// should be expressions
	for _, s := range []string{
		"1 //a",
		"1 /*abc*/",
		"_",
		"3",
		"-3",
		"-a",
		"!a",
		"0",
		"_a",
		"a",
		"a + 3",
		"print(3)",
		"print(3, 4)",
		"print()",
		"a",
		"a+3+4",
		"a * 3",
		"(a)",
		"((a))",
		"(a+3)*4",
		"4 * (a + 3)",
		"a == 4",
		"a > 5",
		"a < 6",
		"a >= 5",
		"a <= 6",
		"a != 7",
		"(void)()",
		"'x'",
		"'\\n'",
		"'\\''",
		`"hello"`,
		"`hello`",
		"`hello\nhi`",
		"*x",
		"**x",
		"&x",
		"a[3]",
		"a[b]",
		"a()[3](b)[7+4]",
		"a[:]",
		"a[:3]",
		"a[3:]",
		"a[0:7+4]",
		"[]int{3, 4, 5, 6}",
		"[]int{3, 4, 5, 6,}",
		"[]uint{}",
	} {
		buf := strings.NewReader(s)
		stmts, es := Stmts("test.g", buf)
		if es != nil {
			t.Log(s)
			for _, e := range es {
				t.Log(e)
			}
			t.Fail()
		} else if len(stmts) != 1 {
			t.Log(s)
			t.Error("should be one statement")
		} else {
			s := stmts[0]
			if _, ok := s.(*ast.ExprStmt); !ok {
				t.Log(s)
				t.Error("should be an expression")
			}
		}
	}

	// should be a statement
	for _, s := range []string{
		";",
		"{;;;;}",
		"{}",
		"{};",
		"{;}",
		"{3}",
		"a = 3",
		"a, b = 4, x",
		"a(), b() = 4, x(x())",
		"a := 3",
		"a := 3+4",
		"a, b := 4, x",
		"for {}",
		"for true { }",
		"for (true) { }",
		"for a == 3 { }",
		"for ;; { }",
		"for ;false; { }",
		"for i:=0; i<3; i++ { }",
		"for ; i<3; i++ { }",
		"for ; i<3; { }",
		"if true { }",
		"if (true) { }",
		"if true { } else { }",
		`if true {
			print(3)
			print(5)
		} else {
			print(4)
		}`,
		"break",
		"continue",
		"if true return",
		"if true break",
		"if true continue",
		"if true return 3",
		"if true { return }",
		"if true { return; break }",
		`for true {
			print(3)
			read()
		}`,
		"var a int",
		"var a int = 3",
		"var a = 3",
		"var a, b int = 3",
		"var a, b int",
		"var a, b int = 3, 4",
		"var ()",
		"var (a, b int)",
		"var (a, b int = 3, 4)",
		"var (a, b = 3, 4)",
		"var (a int; b int)",
		"var (a int\n b int)",
		"var (\n a int \n);",
		"a++",
		"a--",
		"ret := (a.b & 0x1) > 0",
		"a := []int{3,4}",
		"switch 0 { }",
		"switch 0 { case 3: }",
		`switch 0 {
			case 3:
			case 4:
		}`,
		"switch 0 { default: }",
	} {
		buf := strings.NewReader(s)
		stmts, es := Stmts("test.g", buf)
		if es != nil {
			t.Log(s)
			for _, e := range es {
				t.Log(e)
			}
			t.Fail()
		} else if len(stmts) != 1 {
			t.Log(s)
			t.Errorf("should be one statement, got %d", len(stmts))
		} else {
			s := stmts[0]
			if _, ok := s.(*ast.ExprStmt); ok {
				t.Log(s)
				t.Error("should not be an expression")
			}
		}
	}
}

func TestStmts_bad(t *testing.T) {
	o := func(code, input string) {
		buf := strings.NewReader(input)
		stmts, es := Stmts("test.g", buf)
		if es == nil || stmts != nil {
			t.Log(input)
			t.Errorf("should error: %s", code)
			return
		}
		errNum := len(es)
		if errNum != 1 {
			t.Log(input)
			t.Logf("%d errors returned", errNum)
			for _, err := range es {
				t.Log(err.Code)
			}
		}
		if len(code) < 7 || code[:7] != "lexing." {
			code = "pl." + code
		}
		if es[0].Code != code {
			t.Log(input)
			t.Log(es[0].Err)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
			return
		}
	}

	o("missingSemi", "3 3")
	o("missingSemi", "3a")
	o("missingSemi", "3x3")
	o("missingSemi", "if true break else continue")
	o("missingSemi", "if true break else {}")
	o("missingSemi", "if true break return")
	o("missingSemi", "if true { x{ } else {}")
	o("missingSemi", "if true { { } else {}")
	o("missingSemi", "if true { x{ } else {}")
	o("missingSemi", "var a b c")
	o("missingSemi", "var a b c = 3, 4")
	o("missingSemi", "var a b = c d")
	o("missingSemi", "a]")

	o("missingSemi", "var (a int, b int);")

	o("expectOperand", "p(")
	o("expectOperand", "{}}")
	o("expectOperand", "p(;)")
	o("expectOperand", "p())")
	o("expectOperand", "}")
	o("expectOperand", "()")
	o("expectOperand", "if { }")
	o("expectOperand", "if { else }")
	o("expectOperand", "a, b := ()")
	o("expectOperand", "for ; {}")
	o("expectOperand", "for ; ")
	o("expectOperand", "for true ;")
	o("expectOperand", "if true { x( } else {}")
	o("expectOperand", "if true { x(;) } else {}")
	o("expectOperand", "if true { x( } else {}")
	o("expectOperand", "if true { x( } else {}")
	o("expectOperand", "a[]")
	o("expectOperand", "a[")
	o("expectOperand", "++i")
	o("expectOperand", "a := []int{,}")

	o("expectOp", "{")
	o("expectOp", "if true { ")
	o("expectOp", "for ;;;; {}")
	o("expectOp", "for ;;; {}")
	o("expectOp", "var ( a int;")
	o("expectOp", "(i++)+3")
	o("expectOp", "a := []int{3, 4, 5, 6")
	o("expectOp", "a := []int{3\n}")

	o("expectType", "var a")
	o("expectType", "var (a)")
	o("expectType", "var a)")

	o("missingIfBody", "if true; { }")
	o("missingIfBody", "if true _:=true")
	o("missingIfBody", "if true else {}")
	o("missingIfBody", "if true {} else; { }")
	o("elseStart", "if true { }; else {}")
	o("elseStart", "if true break; else {}")

	// switch
	o("missingSwitch", "case")
	o("missingSwitch", "default")
	o("missingCaseInSwitch", `switch 1 {b:=2}`)
	o("invalidFallthrough", `switch 2 { case 2: fallthrough}`)
	o("invalidFallthrough", `switch 2 { case 1:
		fallthrough;fallthrough;case 2:}`)
	o("invalidFallthrough", `switch 2 { case 2:
		if true {fallthrough}}`)
	o("invalidFallthrough", "fallthrough")

	o("lexing.unexpected", "var = 3")
	o("lexing.unexpected", "var \n ()")
	o("lexing.unexpected", "var {}")

	o("lexing.unexpectedEOF", "`")
	o("lexing.unexpectedEOF", "`x")
	o("lexing.unexpectedEOF", "/*")
	o("lexing.unexpectedEOF", "\"")
	o("lexing.unexpectedEndl", "\"\n")

	o("lexing.unknownESC", "'\\\"'")

	o("illegalChar", "@")
	o("invalidDotDot", "..")
	o("incOnExprList", "a,b++")

}
