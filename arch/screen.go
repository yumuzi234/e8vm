package arch

import (
	"e8vm.io/e8vm/arch/vpc"
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

const (
	screenWidth  = 80
	screenHeight = 24
)

func (s *screen) Handle(req []byte) ([]byte, int32) {
	dec := coder.NewDecoder(req)
	cmd := dec.U8()
	if dec.Err != nil {
		return nil, vpc.ErrInvalidArg
	}

	switch cmd {
	case 0, 1:
		c := dec.U8()
		line := uint32(dec.U8())
		col := uint32(dec.U8())

		if dec.Err != nil {
			return nil, vpc.ErrInvalidArg
		}

		if cmd == 0 {
			s.textUpdate[line*screenWidth+col] = c
		} else { // cmd == 1
			s.colorUpdate[line*screenWidth+col] = c
		}
	default:
		return nil, vpc.ErrInvalidArg
	}

	return nil, 0
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
