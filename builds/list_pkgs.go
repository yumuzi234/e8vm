package builds

import (
	"fmt"
	"path"
	"sort"
)

// ListPkgs lists all packages based on the selector.
// If a selector ends with "...", it selects all sub packages.
// Otherwise, it selects the exact package.
func ListPkgs(input Input, selectors []string) ([]string, error) {
	picked := make(map[string]struct{})
	add := func(ps []string) {
		for _, p := range ps {
			if !IsPkgPath(p) {
				continue
			}
			picked[p] = struct{}{}
		}
	}

	for _, s := range selectors {
		if s == "*" {
			add(input.Pkgs(""))
			continue
		}

		base := path.Base(s)
		if base == "..." || base == "*" {
			pre := path.Dir(s)
			if pre == "." {
				pre = ""
			}
			pkgs := input.Pkgs(pre)
			if len(pkgs) == 0 {
				err := fmt.Errorf("%q matches no package", s)
				return nil, err
			}
			add(pkgs)
			continue
		}

		if !IsPkgPath(s) {
			err := fmt.Errorf("%q is not a valid package path", s)
			return nil, err
		}
		if !input.HasPkg(s) {
			err := fmt.Errorf("package %q not found", s)
			return nil, err
		}

		add([]string{s})
	}

	if len(picked) == 0 {
		return nil, nil
	}

	ret := make([]string, 0, len(picked))
	for p := range picked {
		ret = append(ret, p)
	}
	sort.Strings(ret)
	return ret, nil
}
