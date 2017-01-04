package sempass

import (
	"shanhu.io/smlvm/lexing"
)

type topoNode struct {
	name string
	tok  *lexing.Token
	deps []string

	queuing bool
	queued  bool
}

type topoSorter struct {
	typ     string
	errCode string
	m       map[string]*topoNode
	order   []string
	err     bool
}

func newTopoSorter(t, errCode string) *topoSorter {
	return &topoSorter{
		typ: t, errCode: errCode, m: make(map[string]*topoNode),
	}
}

func (s *topoSorter) addNode(name string, tok *lexing.Token, deps []string) {
	if _, found := s.m[name]; found {
		panic("duplicate node")
	}

	s.m[name] = &topoNode{name: name, tok: tok, deps: deps}
}

func (s *topoSorter) push(b *builder, node *topoNode) {
	node.queuing = true
	name := node.name

	for _, dep := range node.deps {
		if dep == name {
			b.CodeErrorf(
				node.tok.Pos, s.errCode,
				"%s %s depends on itself", s.typ, name,
			)
			s.err = true
			continue
		}

		depNode := s.m[dep]
		if depNode == nil || depNode.queued {
			continue
		}

		if depNode.queuing {
			b.CodeErrorf(
				node.tok.Pos, s.errCode,
				"%s %s circular depends on %s %s",
				s.typ, node.tok.Lit, s.typ, dep,
			)
			s.err = true
			continue
		}

		s.push(b, depNode)
	}

	node.queuing = false
	node.queued = true
	s.order = append(s.order, node.name)
}

func (s *topoSorter) sort(b *builder) []string {
	s.order = nil
	s.err = false
	for _, node := range s.m {
		if !node.queued {
			s.push(b, node)
		}
	}

	s.m = make(map[string]*topoNode) // clear the map

	if s.err {
		return nil
	}
	return s.order
}
