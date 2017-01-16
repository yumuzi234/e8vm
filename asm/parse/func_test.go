package parse

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"shanhu.io/smlvm/asm/ast"
	"shanhu.io/smlvm/fmtutil"
	"shanhu.io/smlvm/lexing"
)

func pf(s string) (string, []*lexing.Error) {
	r := strings.NewReader(s)
	rc := ioutil.NopCloser(r)
	out := new(bytes.Buffer)
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
		return "", errs
	}

	for _, f := range fs {
		fmt.Fprintf(out, "func %s {\n", f.Name.Lit)
		for _, stmt := range f.Stmts {
			for i, op := range stmt.Ops {
				if i == 0 {
					fmt.Fprint(out, "    ")
				} else {
					fmt.Fprint(out, " ")
				}

				fmt.Fprint(out, op.Lit)
			}
			fmt.Fprintln(out)
		}

		fmt.Fprintf(out, "}\n")
	}
	return out.String(), nil
}

func TestParseFunc(t *testing.T) {
	o := func(in, want string) {
		got, errs := pf(in)
		if errs != nil {
			for _, err := range errs {
				t.Log(err)
			}
			t.Fail()
			return
		}

		want = strings.TrimSpace(fmtutil.BoxSpaceIndent(want))
		got = strings.TrimSpace(got)
		if got != want {
			t.Errorf("parse %q, got %q, want %q", in, got, want)
		}
	}

	o(`
	func main {
		add r4 /*inline comment*/ r3 r5

		// blank lines are ignored
		sub r0   r0		r1 // some comment
		/* some block comment also */
	}`, `
	func main {
		add r4 r3 r5
		sub r0 r0 r1
	}`)

	o("func main {}", `
	func main {
	}`)

	o(`
	func main {
	}`, `
	func main {
	}`)
}
