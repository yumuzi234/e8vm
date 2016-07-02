package arch8

type screen struct {
	ptext  *page
	pcolor *page
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
	// check if has incoming pulling
	// and returns update if possible.
}
