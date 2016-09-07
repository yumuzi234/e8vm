package arch

import (
	"e8vm.io/e8vm/coder"
)

// Screen is an interface for drawing characters.
type Screen interface {
	NeedUpdate() bool
	UpdateText(m map[uint32]byte)
	UpdateColor(m map[uint32]byte)
}

type screen struct {
	textUpdate  map[uint32]byte
	colorUpdate map[uint32]byte
	s           Screen
}

func newScreen(s Screen) *screen {
	if s == nil {
		panic("creating nil screen")
	}

	return &screen{
		textUpdate:  make(map[uint32]byte),
		colorUpdate: make(map[uint32]byte),
		s:           s,
	}
}

func (s *screen) Handle(req, resp []byte) (n, res uint32) {
	dec := coder.NewDecoder(req)
	cmd := dec.U8()
	if dec.Err != nil {
		return 0, callsResInvalidRequest
	}

	const screenWidth = 80

	switch cmd {
	case 0, 1:
		c := dec.U8()
		line := uint32(dec.U8())
		col := uint32(dec.U8())

		if dec.Err != nil {
			return 0, callsResInvalidRequest
		}

		if cmd == 0 {
			s.textUpdate[line*screenWidth+col] = c
		} else { // cmd == 1
			s.colorUpdate[line*screenWidth+col] = c
		}
	default:
		return 0, callsResInvalidRequest
	}

	return 0, 0
}

func (s *screen) flush() {
	if len(s.textUpdate) > 0 {
		s.s.UpdateText(s.textUpdate)
		s.textUpdate = make(map[uint32]byte)
	}
	if len(s.colorUpdate) > 0 {
		s.s.UpdateColor(s.colorUpdate)
		s.colorUpdate = make(map[uint32]byte)
	}
}

func (s *screen) Tick() {
	if s.s.NeedUpdate() {
		s.flush()
	}
}
