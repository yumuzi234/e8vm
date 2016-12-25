package devs

// Keyboard manages keyboard inputs.
type Keyboard struct {
	sender Sender
}

// NewKeyboard creates a new keyboard handler.
func NewKeyboard(s Sender) *Keyboard {
	return &Keyboard{sender: s}
}

// KeyDown sends in a key down event.
func (k *Keyboard) KeyDown(code uint8) error {
	k.sender.Send([]byte{0, code})
	return nil
}
