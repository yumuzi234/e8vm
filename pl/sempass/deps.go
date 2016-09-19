package sempass

import (
	"path/filepath"

	"shanhu.io/smlvm/dagvis"
	"shanhu.io/smlvm/pl/ast"
)

type deps map[string]map[string]struct{}

func newDeps(files map[string]*ast.File) deps {
	ret := make(deps)

	for _, f := range files {
		// use the path in ast, this is consistent with
		// the pos filename.
		ret[f.Path] = make(map[string]struct{})
	}

	return ret
}

func (d deps) add(refFile, defFile string) bool {
	fileDeps, ok := d[refFile]
	if !ok {
		return false
	}
	_, ok = d[defFile]
	if !ok {
		return false
	}

	fileDeps[defFile] = struct{}{}
	return true
}

func (d deps) graph() *dagvis.Graph {
	ret := make(map[string][]string)
	bases := make(map[string]struct{})
	nameMap := make(map[string]string)

	for f := range d {
		base := filepath.Base(f)
		if _, ok := bases[base]; ok {
			panic("dup file")
		}
		bases[base] = struct{}{}
		nameMap[f] = base
	}

	for f, deps := range d {
		base := nameMap[f]
		depList := make([]string, 0, len(deps))
		for dep := range deps {
			depBase, ok := nameMap[dep]
			if !ok {
				panic("missing file")
			}
			depList = append(depList, depBase)
		}
		ret[base] = depList
	}

	g := &dagvis.Graph{Nodes: ret}
	return g
}
