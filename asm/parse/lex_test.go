package parse

import (
	"fmt"
	"io/ioutil"
	"strings"

	"shanhu.io/smlvm/lexing"
)

func tstr(t int) string {
	switch t {
	case lexing.EOF:
		return "eof"
	case lexing.Comment:
		return "cm"
	case Keyword:
		return "kw"
	case Operand:
		return "op"
	case String:
		return "str"
	case Lbrace:
		return "lb"
	case Rbrace:
		return "rb"
	case Endl:
		return "endl"
	case lexing.Illegal:
		return "!"
	}
	return fmt.Sprintf("!%d", t)
}

func o(s string) {
	f := "t.s8"
	r := strings.NewReader(s)
	rc := ioutil.NopCloser(r)
	toks, errs := Tokens(f, rc)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		fmt.Printf("%d error(s)\n", len(errs))
	} else {
		for _, t := range toks {
			if t.Type == lexing.EOF || t.Type == Endl ||
				t.Type == Lbrace || t.Type == Rbrace {
				fmt.Printf("%s:%d: %s\n", t.Pos.File, t.Pos.Line,
					tstr(t.Type))
			} else {
				fmt.Printf("%s:%d: %s - %q\n", t.Pos.File, t.Pos.Line,
					tstr(t.Type), t.Lit)
			}
		}
	}
}

func ExampleTokens_case1() {
	o("\n")
	// Output:
	// t.s8:1: endl
	// t.s8:1: eof
}

func ExampleTokens_case2() {
	o("")
	// Output:
	// t.s8:1: eof
}

func ExampleTokens_case3() {
	o("func a { // comment \n\tsyscall\n}")
	// Output:
	// t.s8:1: kw - "func"
	// t.s8:1: op - "a"
	// t.s8:1: lb
	// t.s8:1: cm - "// comment "
	// t.s8:1: endl
	// t.s8:2: op - "syscall"
	// t.s8:2: endl
	// t.s8:3: rb
	// t.s8:3: eof
}

func ExampleTokens_case4() {
	o("func a{}")
	// Output
	// t.s8:1: kw - "func"
	// t.s8:1: op - "a"
	// t.s8:1: lb
	// t.s8:1: rb
}

func ExampleTokens_keywords() {
	o("func var const import")
	// Output:
	// t.s8:1: kw - "func"
	// t.s8:1: kw - "var"
	// t.s8:1: kw - "const"
	// t.s8:1: kw - "import"
	// t.s8:1: eof
}

func ExampleTokens_string() {
	o("\"some string \\\"\\\\ here\"")
	// Output:
	// t.s8:1: str - "\"some string \\\"\\\\ here\""
	// t.s8:1: eof
}

func ExampleTokens_badstr1() {
	o(`"some string`)
	// Output:
	// t.s8:1: unexpected eof in string
	// 1 error(s)
}

func ExampleTokens_badstr2() {
	o("\"some string\n")
	// Output:
	// t.s8:1: unexpected endl in string
	// 1 error(s)
}

func ExampleTokens_badcomment() {
	o(`/*some comment`)
	// Output:
	// t.s8:1: unexpected eof in block comment
	// 1 error(s)
}

func ExampleTokens_comments() {
	o("// line comment \n /* block comment */")
	// Output:
	// t.s8:1: cm - "// line comment "
	// t.s8:1: endl
	// t.s8:2: cm - "/* block comment */"
	// t.s8:2: eof
}
