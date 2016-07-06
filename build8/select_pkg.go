package build8

import (
	"fmt"
	"strings"
)

// SelectPkgs selects the package to build based on the selector.
func SelectPkgs(in Input, s string) ([]string, error) {
	if s == "" {
		return in.Pkgs(""), nil
	}
	if strings.HasSuffix(s, "...") {
		prefix := strings.TrimSuffix(s, "...")
		return in.Pkgs(prefix), nil
	}

	if !in.HasPkg(s) {
		return nil, fmt.Errorf("package %q not found", s)
	}
	return []string{s}, nil
}
