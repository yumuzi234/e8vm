package gfmt

import (
	"bytes"
	"strings"
	"testing"

	"e8vm.io/e8vm/g8/parse"
)

func TestFormatFile(t *testing.T) {
	gfmt := func(s string) string {
		r := strings.NewReader(s)
		ast, rec, es := parse.File("", r, false)
		if len(es) > 0 {
			t.Errorf("parsing %q failed with errors", s)
			for _, err := range es {
				t.Log(err)
			}
		}

		out := new(bytes.Buffer)
		FprintFile(out, ast, rec)
		return out.String()
	}

	o := func(s, exp string) {
		got := gfmt(s)
		if exp != got {
			t.Errorf("gfmt %q: expect %q, got %q", s, exp, got)
		}
	}

	o("func main() {}", "func main() {}\n")         // tailing end line
	o("func main  () {  }", "func main() {}\n")     // remove spaces
	o("\n\nfunc main  () {  }", "func main() {}\n") // remove lines
	o("func main(){}", "func main() {}\n")          // add spaces
	o("func main() {\n}", "func main() {}\n")       // merge oneliner

	/*
		TODO(kcnm): pass this
		o("func main() { var a [5]int; b := a[:] }",
			"func main() {\n    var a [5]int\n    b := a[:]\n}\n")
	*/

	o("// some comment", "// some comment\n") // comment
}
