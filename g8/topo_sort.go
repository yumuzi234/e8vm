package g8

import (
	"e8vm.io/e8vm/lex8"
)

type topoSortNode struct {
	name string
	tok  *lex8.Token
	deps []string

	queuing bool
	queued  bool
}

type topoSorter struct {
	typ   string
	m     map[string]*topoSortNode
	order []*topoSortNode
}

func newTopoSorter(t string) *topoSorter {
	return &topoSorter{typ: t, m: make(map[string]*topoSortNode)}
}

func (s *topoSorter) push(log lex8.Logger, node *topoSortNode) {
	node.queuing = true
	name := node.name

	for _, dep := range node.deps {
		if dep == name {
			log.Errorf(node.tok.Pos, "%s %s depends on itself",
				s.typ, name,
			)
			continue
		}

		depNode := s.m[dep]
		if depNode == nil || depNode.queued {
			continue
		}

		if depNode.queuing {
			log.Errorf(node.tok.Pos,
				"%s %s circular depends on %s %s",
				s.typ, node.tok.Lit, s.typ, dep,
			)
			continue
		}

		s.push(log, depNode)
	}

	node.queuing = false
	node.queued = true
	s.order = append(s.order, node)
}

func (s *topoSorter) sort(log lex8.Logger) []*topoSortNode {
	s.order = nil
	for _, node := range s.m {
		if !node.queued {
			s.push(log, node)
		}
	}
	return s.order
}

func sortStructs(b *builder, m map[string]*structInfo) []*structInfo {
	s := newTopoSorter("struct")

	for name, info := range m {
		s.m[info.name.Lit] = &topoSortNode{
			name: name,
			tok:  info.name,
			deps: info.deps,
		}
	}

	order := s.sort(b)
	var ret []*structInfo
	for _, node := range order {
		ret = append(ret, m[node.name])
	}

	return ret
}

func sortConsts(b *builder, m map[string]*constInfo) []*constInfo {
	s := newTopoSorter("const")

	for name, info := range m {
		s.m[info.name.Lit] = &topoSortNode{
			name: name,
			tok:  info.name,
			deps: info.deps,
		}
	}

	order := s.sort(b)
	var ret []*constInfo
	for _, node := range order {
		ret = append(ret, m[node.name])
	}

	return ret
}
