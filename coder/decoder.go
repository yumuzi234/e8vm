package coder

import (
	"bytes"
	"encoding/binary"
)

// Decoder is a simple binary decoder
type Decoder struct {
	r   *bytes.Reader
	Err error
}

// NewDecoder creates a new decoder around the buffer.
func NewDecoder(buf []byte) *Decoder {
	return &Decoder{
		r: bytes.NewReader(buf),
	}
}

// U8 reads a byte out of the decoder.
func (c *Decoder) U8() byte {
	ret, err := c.r.ReadByte()
	if err != nil {
		c.Err = err
		return 0
	}
	return ret
}

// U32 reads a word out of the decoder.
func (c *Decoder) U32() uint32 {
	var buf [4]byte
	if _, err := c.r.Read(buf[:]); err != nil {
		c.Err = err
		return 0
	}
	return binary.LittleEndian.Uint32(buf[:])
}
