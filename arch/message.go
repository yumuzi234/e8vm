package arch

// Message is a generic message among different components.
type Message struct {
	From    uint8
	To      uint8
	Payload []byte
}
