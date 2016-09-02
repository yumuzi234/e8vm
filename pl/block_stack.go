package pl

import (
	"e8vm.io/e8vm/pl/codegen"
)

type blockLayer struct {
	name string
	b    *codegen.Block
}

type blockStack struct {
	bs   []*blockLayer
	bmap map[string]*blockLayer
}

func newBlockStack() *blockStack {
	ret := new(blockStack)
	ret.bmap = make(map[string]*blockLayer)
	return ret
}

func (s *blockStack) push(b *codegen.Block, name string) bool {
	if name != "" && s.bmap[name] != nil {
		return false
	}

	layer := &blockLayer{name, b}
	s.bs = append(s.bs, layer)
	if name != "" {
		s.bmap[name] = layer
	}
	return true
}

func (s *blockStack) pop() {
	nlayer := len(s.bs)
	poped := s.bs[nlayer-1]
	s.bs = s.bs[:nlayer-1]
	if poped.name != "" {
		delete(s.bmap, poped.name)
	}
}

func (s *blockStack) top() *codegen.Block {
	nlayer := len(s.bs)
	if nlayer == 0 {
		return nil
	}
	ret := s.bs[nlayer-1]
	return ret.b
}

func (s *blockStack) byName(name string) *codegen.Block {
	layer := s.bmap[name]
	if layer == nil {
		return nil
	}
	return layer.b
}
