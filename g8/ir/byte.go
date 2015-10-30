package ir

type byt struct{ v uint8 } // true or false

func (b *byt) String() string {
	if b.v == 0 {
		return "false"
	}
	return "true"
}

func (b *byt) Size() int32 { return 1 }

func (b *byt) RegSizeAlign() bool { return false }

// Byt creates a new byte referece of a byte
func Byt(b uint8) Ref { return &byt{b} }
