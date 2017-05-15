package net

// Handler handles an incoming packet.
type Handler interface {
	HandlePacket(p []byte) error
}
