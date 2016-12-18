package devs

// Sender sends a packet.
type Sender interface {
	Send(p []byte)
}
