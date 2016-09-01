package builds

import (
	"e8vm.io/e8vm/lexing"
)

func prepare(c *context, p string) (*pkg, []*lexing.Error) {
	saved := c.pkgs[p]
	if saved != nil {
		return saved, nil // already prepared
	}

	pkg := newPkg(c.input, c.output, p)
	c.pkgs[p] = pkg
	if pkg.err != nil {
		return pkg, nil
	}

	if es := pkg.lang.Prepare(pkg.srcMap(), pkg); es != nil {
		return pkg, es
	}

	// recursively prepare imported packages
	for _, imp := range pkg.imports {
		impPkg, es := prepare(c, imp.Path)
		if es != nil {
			return pkg, es
		}

		if impPkg.err != nil {
			return pkg, []*lexing.Error{{
				Pos: imp.Pos,
				Err: impPkg.err,
			}}
		}
	}

	// mount the deps
	deps := make([]string, 0, len(pkg.imports))
	for _, imp := range pkg.imports {
		deps = append(deps, imp.Path)
	}
	c.deps[p] = deps

	return pkg, nil
}
