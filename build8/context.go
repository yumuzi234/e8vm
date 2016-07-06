package build8

import (
	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/link8"
)

type context struct {
	input  Input
	output Output
	*Options

	pkgs map[string]*pkg
	deps map[string][]string

	linkPkgs   map[string]*link8.Pkg
	debugFuncs *debug8.Funcs
}
