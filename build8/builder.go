package build8

import (
	"fmt"

	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// Builder builds a bunch of packages.
type Builder struct {
	*context
}

// NewBuilder creates a new builder with a particular home directory
func NewBuilder(input Input, output Output) *Builder {
	return &Builder{
		context: &context{
			input:      input,
			output:     output,
			pkgs:       make(map[string]*pkg),
			deps:       make(map[string][]string),
			linkPkgs:   make(map[string]*link8.Pkg),
			debugFuncs: debug8.NewFuncs(),
			Options:    new(Options),
		},
	}
}

// BuildPkgs builds a list of packages
func (b *Builder) BuildPkgs(pkgs []string) []*lex8.Error {
	return build(b.context, pkgs)
}

// Build builds a package.
func (b *Builder) Build(p string) []*lex8.Error {
	if !b.input.HasPkg(p) {
		return lex8.SingleErr(fmt.Errorf(
			"package %q not found", p,
		))
	}
	return b.BuildPkgs([]string{p})
}

// BuildPrefix builds packages with a particular prefix.
// in the path.
func (b *Builder) BuildPrefix(repo string) []*lex8.Error {
	return b.BuildPkgs(b.input.Pkgs(repo))
}

// BuildAll builds all packages.
func (b *Builder) BuildAll() []*lex8.Error { return b.BuildPrefix("") }

// Plan returns all the packages required for building the specified
// target packages.
func (b *Builder) Plan(pkgs []string) ([]string, []*lex8.Error) {
	for _, p := range pkgs {
		if _, es := prepare(b.context, p); es != nil {
			return nil, es
		}
	}

	g := &dagvis.Graph{b.deps}
	g = g.Reverse()

	ret, err := dagvis.TopoSort(g)
	if err != nil {
		return nil, lex8.SingleErr(err)
	}
	return ret, nil
}
