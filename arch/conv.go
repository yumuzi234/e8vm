// Package arch8 emulates the E8VM's instruction set.
package arch

import (
	"encoding/binary"
)

const (
	// InitPC points the default starting program counter
	InitPC = 0x8000
)

// Endian is the machine's endian (byte order).
var Endian = binary.LittleEndian
