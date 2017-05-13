package net

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	headerLen      = 16
	lenOffset      = 2
	destIPOffset   = 4
	srcIPOffset    = 8
	destPortOffset = 12
	srcPortOffset  = 14
)

var coding = binary.BigEndian

func checkHeaderLen(p []byte) error {
	if len(p) < headerLen {
		return errors.New("incomplete header")
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

// AddrStr returns the IPv4 representation of the given IP address.
func AddrStr(addr uint32) string {
	b0 := addr & 0xff
	b1 := (addr >> 8) & 0xff
	b2 := (addr >> 16) & 0xff
	b3 := (addr >> 24) & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", b3, b2, b1, b0)
}
