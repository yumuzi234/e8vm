package builds

import (
	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/link"
)

type context struct {
	input  Input
	output Output
	*Options

	pkgs map[string]*pkg
	deps map[string][]string

	linkPkgs   map[string]*link.Pkg
	debugFuncs *debug.Funcs
}
