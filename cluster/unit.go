package cluster

import (
	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/net"
)

// unit is a machine unit in a cluster.
type unit struct {
	m       *arch.Machine
	gateIn  *net.AddrMapper
	gateOut *net.AddrMapper
}

func (u *unit) HandlePacket(p []byte) error {
	// this normally sends the packet to the machine.
	return nil
}
