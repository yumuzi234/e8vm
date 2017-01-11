package gfmt

import (
	"strings"

	"shanhu.io/smlvm/fmtutil"
)

func formatComment(c string) string {
	if strings.HasPrefix(c, "//") {
		// line comment
		if c != "//" && !strings.HasPrefix(c, "// ") {
			c = "// " + strings.TrimPrefix(c, "//")
		}
		c = strings.TrimRight(c, " \t\n")
		return c
	}

	// block comment here
	body := strings.TrimPrefix(c, "/*")
	body = strings.TrimSuffix(body, "*/")
	body = fmtutil.BoxSpaceIndent(body)
	return "/*" + body + "*/"
}
