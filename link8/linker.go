package link8

type linker struct {
	pkgs        []*Pkg
	pkgIndexMap map[string]int
}

func newLinker() *linker {
	ret := new(linker)
	ret.pkgIndexMap = make(map[string]int)
	return ret
}

func (lnk *linker) addPkg(p *Pkg) (index int, isNew bool) {
	path := p.path
	index, found := lnk.pkgIndexMap[path]
	if found {
		return index, false
	}

	index = len(lnk.pkgs)
	lnk.pkgs = append(lnk.pkgs, p)
	lnk.pkgIndexMap[path] = index
	return index, true
}

func (lnk *linker) addPkgs(p *Pkg) int {
	index, isNew := lnk.addPkg(p)
	if isNew {
		for _, req := range p.requires {
			lnk.addPkgs(req)
		}
	}

	return index
}

func (lnk *linker) pkgIndex(path string) int {
	ret, found := lnk.pkgIndexMap[path]
	if !found {
		panic("not found")
	}
	return ret
}

func (lnk *linker) pkg(i int) *Pkg {
	return lnk.pkgs[i]
}

func (lnk *linker) npkg() int {
	return len(lnk.pkgs)
}
