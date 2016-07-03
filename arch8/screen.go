package arch8

// Screen is an interface for drawing characters.
type Screen interface {
	NeedUpdate() bool
	UpdateText(m map[uint32]byte)
	UpdateColor(m map[uint32]byte)
}

type screen struct {
	ptext  *page
	pcolor *page
	s      Screen
}

func newScreen(ptext, pcolor *page, s Screen) *screen {
	ptext.trackDirty()
	pcolor.trackDirty()

	return &screen{
		ptext:  ptext,
		pcolor: pcolor,
		s:      s,
	}
}

func (s *screen) Tick() {
	if s.s == nil { // headless
		return
	}
	if !s.s.NeedUpdate() { // not refreshed yet
		return
	}

	if len(s.ptext.dirty) > 0 {
		s.s.UpdateText(s.ptext.dirtyBytes())
		s.ptext.trackDirty()
	}

	if len(s.pcolor.dirty) > 0 {
		s.s.UpdateColor(s.pcolor.dirtyBytes())
		s.pcolor.trackDirty()
	}
}
