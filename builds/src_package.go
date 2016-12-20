package builds

// SrcPackage contains all the input of a package
type SrcPackage struct {
	Path  string  // package import path
	Lang  string  // language the package is written
	Hash  string  // a signature; empty for unknown
	Files []*File // list of source files
}

// Loader loads packages.
type Loader interface {
	Load(path string) *SrcPackage
}

// Lister lists packages of a given pattern.
type Lister interface {
	List(pattern string) []*SrcPackage
}
