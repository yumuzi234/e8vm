package table

import (
	"shanhu.io/smlvm/arch/vpc"
	"shanhu.io/smlvm/coder"
)

// Table is a virtual card table device.
type Table struct {
	out Render
	in  vpc.Sender
}

// New creates a new virtual card table device.
func New(out Render, in vpc.Sender) *Table {
	return &Table{
		out: out,
		in:  in,
	}
}

const (
	actionNoop = iota
	actionShow
	actionShowFront
	actionShowBack
	actionHide
	actionSetFace
)

var actionStrings = map[uint8]string{
	actionNoop:      "noop",
	actionShow:      "show",
	actionShowFront: "showFront",
	actionShowBack:  "showBack",
	actionHide:      "hide",
	actionSetFace:   "setFace",
}

// Handle handles an incoming VPC.
func (t *Table) Handle(req []byte) ([]byte, int32) {
	if t.out == nil {
		return nil, 0
	}

	dec := coder.NewDecoder(req)
	action := dec.U8()
	pos := dec.U8()
	face := ""
	if action == actionSetFace {
		face = string(rune(dec.U8()))
	}

	if dec.Err != nil {
		return nil, vpc.ErrInvalidArg
	}

	t.out.Act(&Action{
		Action: actionStrings[action],
		Pos:    int(pos),
		Face:   face,
	})

	return nil, 0
}

// Click send in a click.
func (t *Table) Click(pos uint8) error {
	t.in.Send([]byte{pos})
	return nil
}
