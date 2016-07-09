package build8

import (
	"path"
	"sort"
)

// ListPkgs lists all packages based on the selector.
// If a selector ends with "...", it selects all sub packages.
// Otherwise, it selects the exact package.
func ListPkgs(input Input, selectors []string) []string {
	picked := make(map[string]struct{})
	add := func(ps []string) {
		for _, p := range ps {
			picked[p] = struct{}{}
		}
	}

	for _, s := range selectors {
		println(s)
		if s == "*" {
			add(input.Pkgs(""))
		} else if path.Base(s) == "..." {
			pre := path.Dir(s)
			if pre == "." {
				pre = ""
			}
			add(input.Pkgs(pre))
		} else if input.HasPkg(s) {
			println(s)
			add([]string{s})
		}
	}

	if len(picked) == 0 {
		return nil
	}

	ret := make([]string, 0, len(picked))
	for p := range picked {
		ret = append(ret, p)
	}
	sort.Strings(ret)
	return ret
}
