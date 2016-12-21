package builds

import (
	"shanhu.io/smlvm/lexing"
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

	srcPkg := pkg.srcPackage()
	impList, errs := pkg.lang.Prepare(srcPkg.Files)
	if errs != nil {
		return pkg, errs
	}

	// recursively prepare imported packages
	for name, imp := range impList.imps {
		impPath := c.importPath(imp.path)

		pkg.imports[name] = &Import{
			Path: impPath,
			Pos:  imp.pos,
		}

		impPkg, errs := prepare(c, impPath)
		if errs != nil {
			return pkg, errs
		}

		if impPkg.err != nil {
			return pkg, []*lexing.Error{{
				Pos: imp.pos,
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
