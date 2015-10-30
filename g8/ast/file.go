package ast

import (
	"e8vm.io/e8vm/lex8"
)

// PackageTitle is the package clause at the very top.
// It is only required in go-like mode.
type PackageTitle struct {
	Kw   *lex8.Token
	Name *lex8.Token
	Semi *lex8.Token
}

// File is a group of declarations
type File struct {
	Path string // file path

	Package *PackageTitle
	Imports *ImportDecls // optional

	Decls []Decl
}
