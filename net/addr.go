package net

import (
	"fmt"
)

// IPPort is the pair of an IP address and a port.
type IPPort struct {
	IP   uint32
	Port uint16
}

// AddrStr returns the IPv4 representation of the given IP address.
func AddrStr(addr uint32) string {
	b0 := addr & 0xff
	b1 := (addr >> 8) & 0xff
	b2 := (addr >> 16) & 0xff
	b3 := (addr >> 24) & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", b3, b2, b1, b0)
}
