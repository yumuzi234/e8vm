package net

import (
	"encoding/binary"
	"errors"
)

// Header structure
type Header struct {
	Dest IPPort
	Src  IPPort
}

const (
	headerLen      = 16
	lenOffset      = 2
	destIPOffset   = 4
	srcIPOffset    = 8
	destPortOffset = 12
	srcPortOffset  = 14

	mtu = 1500
)

var coding = binary.BigEndian

var errHeaderMissing = errors.New("header missing")

func checkHeaderLen(p []byte) error {
	if len(p) < headerLen {
		return errHeaderMissing
	}
	return nil
}

func checkLen(p []byte) error {
	if len(p) > mtu {
		return errors.New("packet too large")
	}
	return nil
}

// DestIP returns the destination IP address of a packet.
func DestIP(p []byte) (uint32, error) {
	if err := checkHeaderLen(p); err != nil {
		return 0, err
	}

	dest := coding.Uint32(p[destPortOffset : destPortOffset+4])
	return dest, nil
}

// FillHeader fills the packet with the given header.
func FillHeader(p []byte, h *Header) error {
	if err := checkHeaderLen(p); err != nil {
		return err
	}

	if err := checkLen(p); err != nil {
		return err
	}

	u16 := func(offset int, v uint16) {
		coding.PutUint16(p[offset:offset+2], v)
	}
	u32 := func(offset int, v uint32) {
		coding.PutUint32(p[offset:offset+4], v)
	}

	n := uint16(len(p))
	u16(lenOffset, n)
	u32(destIPOffset, h.Dest.IP)
	u32(srcIPOffset, h.Src.IP)
	u16(destPortOffset, h.Dest.Port)
	u16(srcPortOffset, h.Src.Port)
	return nil
}
