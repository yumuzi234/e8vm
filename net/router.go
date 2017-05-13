package net

import (
	"fmt"
)

// Router routes packets to different sinks based on the IP address.
type Router struct {
	hs map[uint32]Handler
}

// NewRouter creates a new router with no routing entries.
func NewRouter() *Router {
	return &Router{
		hs: make(map[uint32]Handler),
	}
}

// SetRoute sets the handler for the given destination address.
func (r *Router) SetRoute(dest uint32, h Handler) {
	if h == nil {
		delete(r.hs, dest)
		return
	}

	r.hs[dest] = h
}

// Handle routes the packet out based on the destination address.
func (r *Router) Handle(p []byte) error {
	dest, err := DestIP(p)
	if err != nil {
		return err
	}

	h, found := r.hs[dest]
	if !found {
		return fmt.Errorf("destination %s not found", AddrStr(dest))
	}
	return h.Handle(p)
}
