package net

// Handler handles an incoming packet.
type Handler interface {
	Handle(p []byte) error
}
