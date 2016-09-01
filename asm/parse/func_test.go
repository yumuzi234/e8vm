package parse

import (
	"fmt"
	"io/ioutil"
	"strings"

	"e8vm.io/e8vm/asm/ast"
	"e8vm.io/e8vm/lexing"
)

func pf(s string) {
	r := strings.NewReader(s)
	rc := ioutil.NopCloser(r)
	p, _ := newParser("t.s8", rc)
	var fs []*ast.Func

	for {
		if p.See(lexing.EOF) {
			break
		}

		f := parseFunc(p)
		if f == nil {
			break
		}

		fs = append(fs, f)
	}

	errs := p.Errs()
	if errs != nil {
		for _, e := range errs {
			fmt.Println(e)
		}
	} else {
		for _, f := range fs {
			fmt.Printf("func %s {\n", f.Name.Lit)
			for _, stmt := range f.Stmts {
				for i, op := range stmt.Ops {
					if i == 0 {
						fmt.Print("    ")
					} else {
						fmt.Print(" ")
					}

					fmt.Print(op.Lit)
				}
				fmt.Println()
			}

			fmt.Printf("}\n")
		}
	}
}

func ExampleFunc_1() {
	pf(`
	func main {
		add r4 /*inline comment*/ r3 r5

		// blank lines are ignored
		sub r0   r0		r1 // some comment
		/* some block comment also */
	}`)
	// Output:
	// func main {
	//     add r4 r3 r5
	//     sub r0 r0 r1
	// }
}

func ExampleFunc_2() {
	pf(`
	func main {}
	`)
	// Output:
	// func main {
	// }
}

func ExampleFunc_3() {
	pf(`
	func main {
	}
	`)
	// Output:
	// func main {
	// }
}

func ExampleFunc_4() {
	pf(`
	func main t {
	}
	`)
	// Output:
	// t.s8:2: expect '{', got operand
}
