package builds

import (
	"strings"
)

func relPath(p string) string {
	return strings.TrimPrefix(p, "/")
}
