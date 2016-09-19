package builds

import (
	"shanhu.io/smlvm/debug"
	link8 "shanhu.io/smlvm/link"
)

type context struct {
	input  Input
	output Output
	*Options

	pkgs map[string]*pkg
	deps map[string][]string

	linkPkgs   map[string]*link8.Pkg
	debugFuncs *debug.Funcs
}
