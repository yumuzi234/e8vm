package ast

import (
	"e8vm.io/e8vm/lexing"
)

// PackageTitle is the package clause at the very top.
// It is only required in go-like mode.
type PackageTitle struct {
	Kw   *lexing.Token
	Name *lexing.Token
	Semi *lexing.Token
}

// File is a group of declarations
type File struct {
	Path string // file path

	Package *PackageTitle
	Imports *ImportDecls // optional

	Decls []Decl
}
