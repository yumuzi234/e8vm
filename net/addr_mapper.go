package net

import (
	"errors"
	"fmt"
)

// AddrMap is an address mapping.
type AddrMap struct {
	M map[uint32]uint32
}

func (m *AddrMap) applyAt(p []byte, offset int) bool {
	bs := p[offset : offset+4]
	ip := coding.Uint32(bs)
	mapped, found := m.M[ip]
	if !found {
		return false
	}
	coding.PutUint32(bs, mapped)
	return true
}

// Apply maps the addresses in the packet to another address set.
func (m *AddrMap) Apply(p []byte) (bool, error) {
	if len(p) < headerLen {
		return false, errHeaderMissing
	}
	found := m.applyAt(p, destIPOffset)
	found = found && m.applyAt(p, srcIPOffset)
	return found, nil
}

// Revert reverts the address mapping.
func (m *AddrMap) Revert() (*AddrMap, error) {
	r := make(map[uint32]uint32)
	for from, to := range m.M {
		if _, found := r[to]; found {
			return nil, fmt.Errorf("not a 1-to-1 mapping")
		}
		r[to] = from
	}
	return &AddrMap{M: r}, nil
}

// AddrMapper is a address mapping filter that translates all
// addresses from one set to another.
type AddrMapper struct {
	Map       *AddrMap
	AllowWild bool
	Out       Handler
}

var errUnknownAddress = errors.New("unknown address")

// HandlePacket maps the address in the packet first and then send it to Out.
func (m *AddrMapper) HandlePacket(p []byte) error {
	ok, err := m.Map.Apply(p)
	if err != nil {
		return err
	}
	if !m.AllowWild && !ok {
		// something un mapaped
		return errUnknownAddress
	}
	if m.Out != nil {
		return m.Out.HandlePacket(p)
	}
	return nil
}
