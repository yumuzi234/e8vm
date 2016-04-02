package gfmt

import (
	"fmt"
	"strings"
)

func indentWithSpace(s string) string {
	ret := ""
	for i, c := range s {
		if c == ' ' {
			ret += " "
		} else if c == '\t' {
			ret += "    "
		} else {
			return ret + s[i:]
		}
	}

	return ret
}

func allSpaces(s string) bool {
	for _, r := range s {
		if r != ' ' {
			return false
		}
	}
	return true
}

func indentCount(s string) int {
	for i, r := range s {
		if r != ' ' {
			return i
		}
	}
	return len(s)
}

func formatComment(c string) string {
	if strings.HasPrefix(c, "//") {
		// line comment
		if c != "//" && !strings.HasPrefix(c, "// ") {
			return "// " + strings.TrimPrefix(c, "//")
		}
		return c
	}

	// block comment here
	fmt.Println(c)
	body := strings.TrimPrefix(c, "/*")
	body = strings.TrimSuffix(body, "*/")

	lines := strings.Split(body, "\n")
	nline := len(lines)

	// convert indent into space, make blank lines empty.
	for i := range lines {
		if i == 0 {
			continue
		}
		line := lines[i]
		line = indentWithSpace(line)
		if allSpaces(line) {
			line = ""
		}
		lines[i] = line
	}

	if nline > 1 {
		minIndentCount := indentCount(lines[1])
		for i, line := range lines {
			if i <= 1 {
				continue
			}
			if line == "" {
				continue
			}

			n := indentCount(line)
			if n < minIndentCount {
				minIndentCount = n
			}
		}

		for i := range lines {
			if i == 0 {
				continue
			}
			if lines[i] == "" {
				continue
			}
			// trim prefix of minIndent
			lines[i] = lines[i][minIndentCount:]
		}
	}

	ret := strings.Join(lines, "\n")
	return "/*" + ret + "*/"
}
