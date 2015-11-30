package g8

import (
	"e8vm.io/e8vm/toposort"
)

func sortStructs(b *builder, m map[string]*structInfo) []*structInfo {
	s := toposort.NewSorter("struct")

	for name, info := range m {
		s.AddNode(name, info.name, info.deps)
	}

	order := s.Sort(b)
	var ret []*structInfo
	for _, name := range order {
		ret = append(ret, m[name])
	}

	return ret
}

func sortConsts(b *builder, m map[string]*constInfo) []*constInfo {
	s := toposort.NewSorter("const")
	for name, info := range m {
		s.AddNode(name, info.name, info.deps)
	}

	order := s.Sort(b)
	var ret []*constInfo
	for _, name := range order {
		ret = append(ret, m[name])
	}

	return ret
}
