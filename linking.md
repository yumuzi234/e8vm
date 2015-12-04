# Linker Design

Unlike in assembly, G language might reference types and hence struct methods defined in packages that are not explicitly imported.

Each package have many faces.
- An assembly package is, well, an assembly package
- A G language package is a G language package, the recognizes assembly package
- And both the asm package and the G language package can be linkable.
- So linkable package should be an interface.
- Each package has two parts, the header part and the content part.
- A package can be saved on disk
- A package can be loaded with only the header or the header and the content.
- Only the header part is enough to support other dependant packages.

// a linking package
type link8.Package interface {
	Link() (map[string]*Var, map[string]*Func)
}

- linker is prepared only when linking starts, which is after all the builds
- all packages that hits in the dependency tree will be added into the linker

type build8.Lang interface {
	IsSrc(filename string) bool
	NewPackage(path, src map[string]*File) (Package, error)
}

type build8.Import struct {
	Path string
	As string // import as
	Pos *lex8.Pos // import position
	Package
}

type build8.Package interface {
	Lang() string // language name
	ListImports() ([]*Import, []*lex8.Error) // list of imports
	Build(imports []*Import) []*lex8.Error
	Linkables() (map[string]*Var, map[string]*Func)
	Main() string
	Tests() (tests string, testMain string)
}

// when building, importing a symbol will look into the package (header) rather than just the linkable. linkable and headers are differet.

// a building package must be first a linking package.

At the very high level:
- we need a long lived linker/builder that is maintained by the builder
- it stores package objects of each module.
- package are indexed uniquely by package paths.
- assembly packages exports consts, vars (data blocks) and funcs (code blocks)
- the linker only care about vars and funcs, for linking.
- g language packages exports consts, types, vars and funcs.
- the linker vars and funcs are common entities.
- linking symbol should be referenced as (pkg name, symbol name) only

so each library should return something like

// this is an inherently inter-twined structure.
// this is why we need interfaces 

struct Pkg {
	Path string
	Vars map[string]*Var
	Func map[string]*Func
}

pkg has var, func
var has pkg, func
func has pkg, var


struct Func {
	Pkg *Pkg
	Name string
	*FuncContent (asm/links)
}

struct Importer {
	Var(path, name string) *Var
	Func(path, name string) *Func
}

interface link8.Lib {
	Compile() (map[string]*Var, map[string]*Func)
}

--

missing return check:
- only check if the function has return values of course.
- and it ends with a statement that is not an ending statement.
- a return is an ending statement
- an infinite loop (with no condition) with no breaks is an ending statement
- an if-else is an ending statement only if both blocks are ending statements


