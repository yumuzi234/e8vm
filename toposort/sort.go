// Package toposort topologically sorts a set of nodes based on their
// dependencies.
package toposort

import (
	"shanhu.io/smlvm/lexing"
)

type node struct {
	name string
	tok  *lexing.Token
	deps []string

	queuing bool
	queued  bool
}

// Sorter sorts a DAG by its dependency
type Sorter struct {
	typ   string
	m     map[string]*node
	order []string
	err   bool
}

// NewSorter creates a new topo sorter
func NewSorter(t string) *Sorter {
	return &Sorter{typ: t, m: make(map[string]*node)}
}

// AddNode addes a new node
func (s *Sorter) AddNode(name string, tok *lexing.Token, deps []string) {
	if _, found := s.m[name]; found {
		panic("duplicate node")
	}

	s.m[name] = &node{name: name, tok: tok, deps: deps}
}

func (s *Sorter) push(log lexing.Logger, node *node) {
	node.queuing = true
	name := node.name

	for _, dep := range node.deps {
		if dep == name {
			log.Errorf(node.tok.Pos, "%s %s depends on itself",
				s.typ, name,
			)
			s.err = true
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
			s.err = true
			continue
		}

		s.push(log, depNode)
	}

	node.queuing = false
	node.queued = true
	s.order = append(s.order, node.name)
}

// Sort sorts the added nodes and returns the node order
func (s *Sorter) Sort(log lexing.Logger) []string {
	s.order = nil
	s.err = false
	for _, node := range s.m {
		if !node.queued {
			s.push(log, node)
		}
	}

	s.m = make(map[string]*node) // clear the map

	if s.err {
		return nil
	}
	return s.order
}
