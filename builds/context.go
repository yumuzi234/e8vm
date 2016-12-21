package builds

import (
	"path"

	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/link"
)

type context struct {
	input   Input
	output  Output
	stdPath string
	*Options

	pkgs map[string]*pkg
	deps map[string][]string

	linkPkgs   map[string]*link.Pkg
	debugFuncs *debug.Funcs
}

func (c *context) importPath(p string) string {
	if c.stdPath == "" {
		return path.Join("/", p)
	}

	if path.IsAbs(p) {
		return p
	}
	return path.Join("/", c.stdPath, p)
}
