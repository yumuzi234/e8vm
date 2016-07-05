package build8

import (
	"io"
	"path"

	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/lex8"
)

func deps(node *dagvis.MapNode) []string {
	depNodes := dagvis.AllInsSorted(node)
	ret := make([]string, 0, len(depNodes))
	for _, dep := range depNodes {
		ret = append(ret, dep.Name)
	}
	return ret
}

func fillImports(c *context, p *pkg) {
	for _, imp := range p.imports {
		imp.Package = c.pkgs[imp.Path].pkg
		if imp.Package == nil {
			panic("bug")
		}
	}
}

func buildMain(c *context, p *pkg) []*lex8.Error {
	lib := p.pkg.Lib
	main := p.pkg.Main

	if main == "" || !lib.HasFunc(main) {
		return nil
	}

	log := lex8.NewErrorList()

	fout := c.output.Bin(p.path)
	lex8.LogError(log, link(c, fout, p, main))
	lex8.LogError(log, fout.Close())

	return log.Errs()
}

func parseOutput(c *context, p string) func(f string, toks []*lex8.Token) {
	if c.SaveFileTokens == nil {
		return nil
	}
	return func(file string, tokens []*lex8.Token) {
		c.SaveFileTokens(path.Join(p, file), tokens)
	}
}

func makePkgInfo(c *context, p *pkg) *PkgInfo {
	return &PkgInfo{
		Path:   p.path,
		Src:    p.srcMap(),
		Import: p.imports,

		Flags: &Flags{StaticOnly: c.StaticOnly},
		Output: func(name string) io.WriteCloser {
			return c.output.Output(p.path, name)
		},
		ParseOutput: parseOutput(c, p.path),
		AddFuncDebug: func(name string, pos *lex8.Pos, frameSize uint32) {
			c.debugFuncs.Add(p.path, name, pos, frameSize)
		},
	}
}

func buildPkg(c *context, pkg *pkg) []*lex8.Error {
	fillImports(c, pkg)

	compiled, es := pkg.lang.Compile(makePkgInfo(c, pkg))
	if es != nil {
		return es
	}
	pkg.pkg = compiled
	c.linkPkgs[pkg.path] = pkg.pkg.Lib // add for linking

	if c.StaticOnly { // static analysis stops here
		return nil
	}

	if es := buildMain(c, pkg); es != nil {
		return es
	}
	if !pkg.runTests { // skip running tests
		return nil
	}

	return runPkgTests(c, pkg)
}

func build(c *context, pkgs []string) []*lex8.Error {
	c.targets = pkgs

	for _, p := range pkgs {
		if _, es := prepare(c, p); es != nil {
			return es
		}
	}

	if c.RunTests {
		for _, p := range pkgs {
			c.pkgs[p].runTests = true
		}
	}

	g := &dagvis.Graph{c.deps}
	g = g.Reverse()

	m, err := dagvis.Layout(g)
	if err != nil {
		return lex8.SingleErr(err)
	}
	if c.SaveDeps != nil {
		c.SaveDeps(m)
	}

	nodes := m.SortedNodes()
	for _, node := range nodes {
		if c.Verbose { // report progress
			logln(c, node.Name)
		}

		pkg := c.pkgs[node.Name]
		pkg.deps = deps(node)
		if es := buildPkg(c, pkg); es != nil {
			return es
		}
	}

	return nil
}
