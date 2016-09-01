package build8

import (
	"fmt"

	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/lexing"
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
func (b *Builder) BuildPkgs(pkgs []string) []*lexing.Error {
	return build(b.context, pkgs)
}

// Build builds a package.
func (b *Builder) Build(p string) []*lexing.Error {
	if !b.input.HasPkg(p) {
		return lexing.SingleErr(fmt.Errorf(
			"package %q not found", p,
		))
	}
	return b.BuildPkgs([]string{p})
}

// BuildPrefix builds packages with a particular prefix.
// in the path.
func (b *Builder) BuildPrefix(prefix string) []*lexing.Error {
	return b.BuildPkgs(b.input.Pkgs(prefix))
}

// BuildAll builds all packages.
func (b *Builder) BuildAll() []*lexing.Error { return b.BuildPrefix("") }

// Plan returns all the packages required for building the specified
// target packages.
func (b *Builder) Plan(pkgs []string) ([]string, []*lexing.Error) {
	for _, p := range pkgs {
		if pkg, es := prepare(b.context, p); es != nil {
			return nil, es
		} else if pkg.err != nil {
			return nil, lexing.SingleErr(pkg.err)
		}
	}

	g := &dagvis.Graph{b.deps}
	g = g.Reverse()

	ret, err := dagvis.TopoSort(g)
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
	return ret, nil
}
