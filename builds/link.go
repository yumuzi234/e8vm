package builds

import (
	"io"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/image"
	"shanhu.io/smlvm/link"
)

func linkPkg(c *context, out io.Writer, p *pkg, main string) error {
	var funcs []*link.PkgSym

	addInit := func(p *pkg) {
		name := p.pkg.Init
		if name != "" && p.pkg.Lib.HasFunc(name) {
			funcs = append(funcs, link.NewPkgSym(p.path, name))
		}
	}

	for _, dep := range p.deps {
		addInit(c.pkgs[dep])
	}
	addInit(p)
	funcs = append(funcs, link.NewPkgSym(p.path, main))

	debugTable := debug.NewTable()
	job := link.NewJob(c.linkPkgs, funcs)
	job.InitPC = c.InitPC
	if job.InitPC == 0 {
		job.InitPC = arch.InitPC
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
