package gfmt

import (
	"testing"

	"strings"

	"e8vm.io/e8vm/fmt8"
)

func formatProg(s string) string {
	s = fmt8.Box(s)
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
		out, errs := Format("a.g", s)
		if len(errs) > 0 {
			t.Errorf("parsing %q failed with errors", s)
			for _, err := range errs {
				t.Log(err)
			}
		}
		return out
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
			t.Errorf("gfmt result %q changed to %q", s, got, got2)
		}
	}

	o("func main() {}", "func main() {}\n")         // tailing end line
	o("func main  () {  }", "func main() {}\n")     // remove spaces
	o("\n\nfunc main  () {  }", "func main() {}\n") // remove lines
	o("func main(){}", "func main() {}\n")          // add spaces
	o("func main() {\n}", "func main() {}\n")       // merge oneliner
	o("// some comment", "// some comment\n")       // comment
	o("/* some comment */", "/* some comment */")   // block comment

	// Common case of line break.
	o("func main() { var a [5]int; b := a[:] }",
		"func main() {\n    var a [5]int\n    b := a[:]\n}\n")
	// Preserves additional/optional line breaks in block.
	o("func main() { var a [5]int;\n\n b := a[:] }",
		"func main() {\n    var a [5]int\n\n    b := a[:]\n}\n")
	// Removes redundant line breaks in block.
	o("func main() { var a [5]int;\n\n b := a[:] }",
		"func main() {\n    var a [5]int\n\n    b := a[:]\n}\n")
	o(`
		func main() { var a [5]int;

			b := a[:]
		}
	`, `
		func main() {
			var a [5]int

			b := a[:]
		}
	`)
	/*
		o(`
			func main() {
			/* some comment /
			}
		`, `
			func main() {
				/* some comment /
			}
		`)
	*/

}
