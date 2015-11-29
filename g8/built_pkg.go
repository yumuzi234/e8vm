package g8

import (
	"e8vm.io/e8vm/link8"
	"e8vm.io/e8vm/sym8"
)

type builtPkg struct {
	p      *pkg
	lib    *link8.Pkg
	isBare bool
}

func (p *builtPkg) Lib() *link8.Pkg { return p.lib }
func (p *builtPkg) Main() string    { return startName }
func (p *builtPkg) Init() string    { return "init" }

func (p *builtPkg) Tests() (map[string]uint32, string) {
	if p.isBare {
		return nil, ""
	}

	tests := make(map[string]uint32)
	for i, name := range p.p.testNames {
		tests[name] = uint32(i)
	}

	return tests, testStartName
}

func (p *builtPkg) Symbols() (string, *sym8.Table) {
	if p.isBare {
		return "g8bare", nil
	}
	return "g8", p.p.tops
}
