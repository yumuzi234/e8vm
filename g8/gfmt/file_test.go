package gfmt

import (
	"strings"
	"testing"
)

func allHasPrefix(lines []string, prefix string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(line, prefix) {
			return false
		}
	}
	return true
}

func tabPrefix(line string) string {
	ret := ""
	for _, c := range line {
		if c == '\t' {
			ret += "\t"
		} else {
			return ret
		}
	}
	return ret
}

func trimBlanks(lines []string) []string {
	// Remove the first blank line
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	nline := len(lines)
	if len(lines) > 0 && strings.TrimSpace(lines[nline-1]) == "" {
		lines = lines[:nline-1]
	}
	return lines
}

func tabPrefixToSpace(line string) string {
	prefix := ""
	for i, c := range line {
		if c == '\t' {
			prefix += "    "
		} else {
			return prefix + line[i:]
		}
	}
	return line
}


func formatProg(s string) string {
	lines := strings.Split(s, "\n")
	lines = trimBlanks(lines)

	if len(lines) == 0 {
		return ""
	}

	first := lines[0]
	prefix := tabPrefix(first)
	if allHasPrefix(lines, prefix) {
		for i := range lines {
			line := lines[i]
			if strings.TrimSpace(line) == "" {
				line = ""
			} else {
				line = strings.TrimPrefix(line, prefix)
				line = tabPrefixToSpace(line)
			}
			lines[i] = line
		}
	}

	ret := strings.Join(lines, "\n")
	if !strings.HasSuffix(ret, "\n") {
		ret += "\n"
	}
	return ret
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
		if exp != got {
			t.Errorf("gfmt %q: expect %q, got %q", s, exp, got)
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
}
