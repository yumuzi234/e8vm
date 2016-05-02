package build8

import (
	"sort"
	"strings"
)

// ListPkgs lists all packages based on the selector.
// If a selector ends with "/...", it selects all sub packages.
// Otherwise, it selects the exact package.
func ListPkgs(input Input, selectors []string) []string {
	picked := make(map[string]struct{})
	add := func(ps []string) {
		for _, p := range ps {
			picked[p] = struct{}{}
		}
	}

	for _, s := range selectors {
		if s == "*" {
			add(input.Pkgs(""))
		} else if strings.HasSuffix(s, "/...") {
			pre := strings.TrimSuffix(s, "/...")
			add(input.Pkgs(pre))
		} else {
			if input.HasPkg(s) {
				add([]string{s})
			}
		}
	}

	ret := make([]string, 0, len(picked))
	for p := range picked {
		ret = append(ret, p)
	}
	sort.Strings(ret)
	return ret
}
