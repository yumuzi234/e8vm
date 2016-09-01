package build8

import (
	"io"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/image"
	"e8vm.io/e8vm/link8"
)

func link(c *context, out io.Writer, p *pkg, main string) error {
	var funcs []*link8.PkgSym

	addInit := func(p *pkg) {
		name := p.pkg.Init
		if name != "" && p.pkg.Lib.HasFunc(name) {
			funcs = append(funcs, &link8.PkgSym{p.path, name})
		}
	}

	for _, dep := range p.deps {
		addInit(c.pkgs[dep])
	}
	addInit(p)
	funcs = append(funcs, &link8.PkgSym{p.path, main})

	debugTable := debug8.NewTable()
	job := link8.NewJob(c.linkPkgs, funcs)
	job.InitPC = c.InitPC
	if job.InitPC == 0 {
		job.InitPC = arch8.InitPC
	}
	job.FuncDebug = func(pkg, name string, addr, size uint32) {
		debugTable.LinkFunc(c.debugFuncs, pkg, name, addr, size)
	}
	secs, err := job.Link()
	if err != nil {
		return err
	}

	debugSec, err := debugSection(debugTable)
	if err != nil {
		return err
	}
	secs = append(secs, debugSec)
	return image.Write(out, secs)
}
