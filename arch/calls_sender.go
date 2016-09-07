package arch

import (
	"container/list"
)

type callsMessage struct {
	service uint32
	p       []byte
}

type callsSender struct {
	service uint32
	queue   *list.List
}

func (s *callsSender) Send(bs []byte) {
	m := &callsMessage{
		service: s.service,
		p:       bs,
	}
	s.queue.PushBack(m)
}
