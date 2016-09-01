package link

// PkgSym is a link to a symbol in a particular package
type PkgSym struct {
	// Pkg is the package path
	Pkg string

	// Sym is the symbol name
	Sym string
}
