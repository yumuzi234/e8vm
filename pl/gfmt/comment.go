package gfmt

import (
	"strings"

	"shanhu.io/smlvm/fmtutil"
)

func formatComment(c string) string {
	if strings.HasPrefix(c, "//") {
		// line comment
		if c != "//" && !strings.HasPrefix(c, "// ") {
			return "// " + strings.TrimPrefix(c, "//")
		}
		return c
	}

	// block comment here
	body := strings.TrimPrefix(c, "/*")
	body = strings.TrimSuffix(body, "*/")
	body = fmtutil.BoxSpaceIndent(body)
	return "/*" + body + "*/"
}
