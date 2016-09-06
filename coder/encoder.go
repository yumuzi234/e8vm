package coder

import (
	"bytes"
	"encoding/binary"
)

// Encoder is a simple binary encoder.
type Encoder struct {
	buf *bytes.Buffer
}

// NewEncoder creates an empty binary encoder.
func NewEncoder() *Encoder {
	return &Encoder{
		buf: new(bytes.Buffer),
	}
}

// U8 appends an uint8 into the buffer.
func (c *Encoder) U8(b byte) {
	c.buf.Write([]byte{b})
}

// U32 appends an uint32 into the buffer.
func (c *Encoder) U32(w uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], w)
	c.buf.Write(buf[:])
}
