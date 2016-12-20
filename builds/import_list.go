package builds

import (
	"shanhu.io/smlvm/lexing"
)

type importPos struct {
	path string
	pos  *lexing.Pos
}

// ImportList saves a list of import declarations.
type ImportList struct {
	imps map[string]*importPos
}

// NewImportList creates an empty import list.
func NewImportList() *ImportList {
	return &ImportList{
		imps: make(map[string]*importPos),
	}
}

// Add adds an import into the import list.
func (lst *ImportList) Add(name, path string, pos *lexing.Pos) {
	lst.imps[name] = &importPos{path, pos}
}
