package vpc

// Sender sends a packet.
type Sender interface {
	Send(p []byte)
}
