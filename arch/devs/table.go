package devs

import (
	"shanhu.io/smlvm/coder"
)

// Table is a virtual card table device.
type Table struct {
	out TableRender
	in  Sender
}

// NewTable creates a new virtual card table device.
func NewTable(out TableRender, in Sender) *Table {
	return &Table{
		out: out,
		in:  in,
	}
}

const (
	noop = iota

	cardShow
	cardShowFront
	cardShowBack
	cardHide
	cardHideFront
	cardHideBack
	cardFace

	labelText

	buttonShow
	buttonHide
	buttonText
)

var actionStrings = map[uint8]string{
	noop: "noop",

	cardShow:      "card.show",
	cardShowFront: "card.showFront",
	cardShowBack:  "card.showBack",
	cardHide:      "card.hide",
	cardHideFront: "card.hideFront",
	cardHideBack:  "card.hideBack",
	cardFace:      "card.face",

	labelText: "label.text",

	buttonShow: "button.show",
	buttonHide: "button.hide",
	buttonText: "button.text",
}

// Handle handles an incoming VPC.
func (t *Table) Handle(req []byte) ([]byte, int32) {
	if t.out == nil {
		return nil, 0
	}

	dec := coder.NewDecoder(req)
	action := dec.U8()
	pos := dec.U8()
	text := ""

	switch action {
	case cardFace:
		text = string(rune(dec.U8()))
	case labelText, buttonText:
		n := dec.U8()
		if n > 0 {
			bs := dec.Bytes(int(n))
			text = string(bs)
		}
	}

	if dec.Err != nil {
		return nil, ErrInvalidArg
	}

	t.out.Act(&TableAction{
		Action: actionStrings[action],
		Pos:    int(pos),
		Text:   text,
	})

	return nil, 0
}

func whatCode(what string) uint8 {
	switch what {
	case "card":
		return 1
	case "button":
		return 2
	}
	return 0
}

// Click sends in a click on the table.
func (t *Table) Click(what string, pos uint8) error {
	t.in.Send([]byte{whatCode(what), pos})
	return nil
}
