package ir

import (
	"fmt"
)

type byt struct {
	v  uint8
	u8 bool
} // true or false

func (b *byt) String() string { return fmt.Sprintf("%d", b.v) }

func (b *byt) Size() int32 { return 1 }

func (b *byt) RegSizeAlign() bool { return false }

// Byt creates a new byte referece of a byte
func Byt(b uint8, u8 bool) Ref { return &byt{b, u8} }
