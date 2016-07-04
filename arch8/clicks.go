package arch8

import (
	"container/list"
	"fmt"
)

type click struct {
	line uint8
	col  uint8
}

type clicks struct {
	q      *list.List
	p      *pageOffset
	intBus intBus
}

const clicksBase = 8

func newClicks(p *page, i intBus) *clicks {
	return &clicks{
		q:      list.New(),
		p:      &pageOffset{p, clicksBase},
		intBus: i,
	}
}

func (c *clicks) addClick(line, col uint8) error {
	if line > 24 {
		return fmt.Errorf("line too big: %d", line)
	}
	if col > 80 {
		return fmt.Errorf("col too big: %d", col)
	}

	if c.q.Len() >= 16 {
		return fmt.Errorf("click event queue full")
	}

	c.q.PushBack(&click{line: line, col: col})
	return nil
}

func (c *clicks) Tick() {
	if c.q.Len() == 0 {
		return
	}

	if c.p.readByte(0) != 0 {
		return
	}

	front := c.q.Front()
	pos := front.Value.(*click)
	c.q.Remove(front)

	c.p.writeByte(0, 1)
	c.p.writeByte(1, 0)
	c.p.writeByte(2, pos.line)
	c.p.writeByte(3, pos.line)
}
