package link8

// addPkgs add a package and recursively adds
// all the packages that this package imported.
func addPkgs(pkgs map[string]*Pkg, p *Pkg) {
	exists := pkgs[p.path]
	if exists != nil {
		if exists != p {
			panic("package path conflict")
		}
		return
	}

	pkgs[p.path] = p
	for _, req := range p.imported {
		addPkgs(pkgs, req)
	}
}
