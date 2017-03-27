package builds

import (
	"fmt"
	"strings"
)

func selectPkgs(src *source, s string) ([]string, error) {
	if s == "" || s == "*" || s == "..." {
		return src.allPkgs("/")
	}

	if strings.HasSuffix(s, "/...") {
		prefix := strings.TrimSuffix(s, "/...")
		return src.allPkgs(prefix)
	}

	ok, err := src.hasPkg(s)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("package %q not found", s)
	}
	return []string{s}, nil
}

// SelectPkgs selects the package to build based on the selector.
func SelectPkgs(in Input, lp *LangPicker, s string) ([]string, error) {
	src := newSource(in, lp)
	return selectPkgs(src, s)
}
