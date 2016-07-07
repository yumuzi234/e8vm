package build8

import (
	"strings"

	"e8vm.io/e8vm/lex8"
)

func isPkgPath(p string) bool {
	p = strings.TrimPrefix(p, "/") // might have a leading slash
	if p == "" {
		return false
	}
	subs := strings.Split(p, "/")
	for _, sub := range subs {
		if !lex8.IsPkgName(sub) {
			return false
		}
	}
	return true
}

// IsSubPkg checks if a package is a subpackage of another package.
func IsSubPkg(p, sub string) bool {
	if p == sub {
		return true
	}
	return strings.HasPrefix(sub, p+"/")
}
