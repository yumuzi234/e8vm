package arch

// Component states.
const (
	StateContinue = iota // continues running
	StateHalt            // stops running
	StateSleep           // pauses running until
)

// MsgReceiver is a message receiver that receives for the next message.
type MsgReceiver interface {
	// Len returns the count of pending messages.
	Len() int

	// Receive receives the next message.
	Receive() (*Message, error)
}

// MsgSender is a message sender that can send out messages.
type MsgSender interface {
	Send(m *Message) error
}

// Component is a component that is linked on a bus
type Component interface {
	Run(q *MsgQueue) int
	State() interface{} // returns a JSON marshalable state.
}
