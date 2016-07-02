package arch8

// Monitor is an interface for drawing characters.
type Monitor interface {
	UpdateText(m map[uint32]byte)
	UpdateColor(m map[uint32]byte)
}

type screen struct {
	ptext  *page
	pcolor *page
	m      Monitor
}

func newScreen(ptext, pcolor *page) *screen {
	ptext.trackDirty()
	pcolor.trackDirty()

	return &screen{
		ptext:  ptext,
		pcolor: pcolor,
	}
}

func (s *screen) Tick() {
	if s.m == nil {
		return
	}

	if len(s.ptext.dirty) > 0 {
		s.m.UpdateText(s.ptext.dirtyBytes())
		s.ptext.trackDirty()
	}

	if len(s.pcolor.dirty) > 0 {
		s.m.UpdateColor(s.pcolor.dirtyBytes())
		s.pcolor.trackDirty()
	}
}
