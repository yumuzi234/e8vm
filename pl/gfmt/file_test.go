package gfmt

import (
	"testing"

	"strings"

	"shanhu.io/smlvm/fmtutil"
)

func formatProg(s string) string {
	s = fmtutil.BoxSpaceIndent(s)
	if strings.HasPrefix(s, "\n") {
		s = s[1:]
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}

func TestFormatFile(t *testing.T) {
	gfmt := func(s string) string {
		out, errs := File("a.g", []byte(s))
		if len(errs) > 0 {
			t.Errorf("parsing %q failed with errors", s)
			for _, err := range errs {
				t.Log(err)
			}
		}
		return string(out)
	}

	o := func(s, exp string) {
		s = formatProg(s)
		exp = formatProg(exp)
		got := gfmt(s)
		got2 := gfmt(s)
		if exp != got {
			t.Errorf("gfmt %q: expect %q, got %q", s, exp, got)
		}
		if got2 != got {
			t.Errorf("gfmt %q: result %q changed to %q", s, got, got2)
		}
	}

	o("func main() {}", "func main() {}\n")         // tailing end line
	o("func main  () {  }", "func main() {}\n")     // remove spaces
	o("\n\nfunc main  () {  }", "func main() {}\n") // remove lines
	o("func main(){}", "func main() {}\n")          // add spaces
	o("func main() {\n}", "func main() {\n}\n")     // do not auto merge one liner
	o("func main() {\n  }", "func main() {\n}\n")   // do not auto merge one liner

	// comment
	o("// some comment", "// some comment\n")
	o("// some comment  ", "// some comment\n")
	o("//    some comment  ", "//    some comment\n")
	o("//some comment  ", "// some comment\n")
	o("//some comment  ", "// some comment\n")

	// block comment
	o("/* some comment */", "/* some comment */")

	// Common case of line break.
	o("func main() { var a [5]int; b := a[:] }",
		"func main() {\n    var a [5]int\n    b := a[:]\n}\n")
	// Preserves additional/optional line breaks in block.
	o("func main() { var a [5]int;\n\n b := a[:] }",
		"func main() {\n    var a [5]int\n\n    b := a[:]\n}\n")
	// Removes redundant line breaks in block.
	o("func main() { var a [5]int;\n\n b := a[:] }",
		"func main() {\n    var a [5]int\n\n    b := a[:]\n}\n")
	o("func main() {/*something*/}", "func main() { /*something*/ }\n")
	o("func main() { /*something*/ }", "func main() { /*something*/ }\n")
	o("func main() { // something\n}", "func main() { // something\n}\n")
	o("func main() {\n// something\n}",
		"func main() {\n    // something\n}\n")
	o(`
		func main() { var a [5]int;

			b := a[:]
		}`, `
		func main() {
			var a [5]int

			b := a[:]
		}
	`)
	o(`
		func main() {
			f(); g()
		}`, `
		func main() {
			f()
			g()
		}
	`)
	o(`
		func main() {
			f()   //with some comment       
		}`, `
		func main() {
			f() // with some comment
		}
	`)
	o(`
		func main() {
			f()

			
			g()
		}`, `
		func main() {
			f()

			g()
		}
	`)
	o(`
		func main() {
			f()
		/* some comment */
		}`, `
		func main() {
			f()
			/* some comment */
		}
	`)
	o(`
		func main() {
			f()
		/* some comment */
		}`, `
		func main() {
			f()
			/* some comment */
		}
	`)
	o(`
        import ( "something" )
        `, `
        import (
            "something"
        )
    `)
	o(`
		func main() {
			f(
			1, // arg1
			)
		}`, `
		func main() {
			f(
				1, // arg1
			)
		}
	`)
	o(`
		func main() {
			var a = []int {
			3,   4, 5  ,
				6,   7,   8,
					}
		}`, `
		func main() {
			var a = []int{
				3, 4, 5,
				6, 7, 8,
			}
		}
	`)

	o(`
		func main() {
			f(
				3, /* arg1 */
				4, /* arg2 */
			)
		}
	`, `
		func main() {
			f(
				3, /* arg1 */
				4, /* arg2 */
			)
		}
	`)

	o(`
		func main() {
			a:=2; switch a {case 1:
			_:=3; default  :   falltrhough; case			 2:
			}}`, `
		func main() {
			a := 2
			switch a {
			case 1:
				_ := 3
			default:
				falltrhough
			case 2:
			}
		}
	`)

	o(`
		func main() {a:=2;     if true {} else {}
		if false {{}} else {  _ := 2;}
    		if a == 2 {}
		}`, `
		func main() {
			a := 2
			if true {
			} else {}
			if false {
				{}
			} else {
				_ := 2
			}
			if a == 2 {}
		}
	`)
	o(`
		func main() { if true return }
	`, `
		func main() {
			if true return
		}
	`)
	o(`
		func main() {}; interface T { f(int,     int) (int,
			int); f()}
	`, `
		func main() {}

		interface T {
			f (int, int) (int, int)
			f ()
		}
	`)
	o(`
		func main() {}; struct s {  
			
			a int
			
			
			
			c d}
	`, `
		func main() {}

		struct s {
			a int
			
			c d
		}
	`)
	o(`
		func main() {}; struct s {  
			}
	`, `
		func main() {}

		struct s {}
	`)
}
