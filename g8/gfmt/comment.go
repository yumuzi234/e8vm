package gfmt

import (
	"fmt"
	"strings"
)

func spacePrefix(s string) string {
	ret := ""
	for _, c := range s {
		if c == ' ' {
			ret += " "
		} else if c == '\t' {
			ret += "    "
		} else {
			break
		}
	}
	return ret
}

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

	for i := range lines {
		if i == 0 {
			continue
		}
		lines[i] = indentWithSpace(lines[i])
	}

	prefix := ""
	if nline > 1 {
		prefix = spacePrefix(lines[1])
	}
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], prefix)
	}

	for i := range lines {
		if i == 0 {
			continue
		}
		if i == nline-1 {
			continue
		}
		if strings.TrimSpace(lines[i]) == "" {
			lines[i] = ""
		}
	}

	ret := strings.Join(lines, "\n")
	return "/*" + ret + "*/"
}
