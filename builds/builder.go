package builds

import (
	"fmt"

	"shanhu.io/smlvm/dagvis"
	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/link"
)

// Builder builds a bunch of packages.
type Builder struct {
	*context
}

// NewBuilder creates a new builder with a particular home directory
func NewBuilder(
	input Input, langPicker *LangSet, std string, output Output,
) *Builder {
	src := newSource(input, langPicker)
	return &Builder{
		context: &context{
			src:        src,
			res:        newResults(output),
			stdPath:    std,
			pkgs:       make(map[string]*pkg),
			deps:       make(map[string][]string),
			linkPkgs:   make(map[string]*link.Pkg),
			debugFuncs: debug.NewFuncs(),
			Options:    new(Options),
		},
	}
}

// SelectPkgs selects the package to build based on the selector.
func (b *Builder) SelectPkgs(s string) ([]string, error) {
	return selectPkgs(b.src, s)
}

// BuildPkgs builds a list of packages
func (b *Builder) BuildPkgs(pkgs []string) []*lexing.Error {
	return build(b.context, pkgs)
}

// Build builds a package.
func (b *Builder) Build(p string) []*lexing.Error {
	ok, err := b.src.hasPkg(p)
	if err != nil {
		return lexing.SingleErr(err)
	}
	if !ok {
		err := fmt.Errorf("package %q not found", p)
		return lexing.SingleErr(err)
	}
	return b.BuildPkgs([]string{p})
}

// BuildPrefix builds packages with a particular prefix.
// in the path.
func (b *Builder) BuildPrefix(prefix string) []*lexing.Error {
	pkgs, err := b.src.allPkgs(prefix)
	if err != nil {
		return lexing.SingleErr(err)
	}
	return b.BuildPkgs(pkgs)
}

// BuildAll builds all packages.
func (b *Builder) BuildAll() []*lexing.Error { return b.BuildPrefix("") }

// Plan returns all the packages required for building the specified
// target packages.
func (b *Builder) Plan(pkgs []string) ([]string, []*lexing.Error) {
	for _, p := range pkgs {
		p = b.context.importPath(p)
		if pkg, es := prepare(b.context, p); es != nil {
			return nil, es
		} else if pkg.err != nil {
			return nil, lexing.SingleErr(pkg.err)
		}
	}

	g := dagvis.NewGraph(b.deps).Reverse()
	ret, err := dagvis.TopoSort(g)
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
	return ret, nil
}
