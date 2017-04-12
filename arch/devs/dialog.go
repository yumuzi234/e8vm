package devs

// Dialog is a chat service.
// TODO(h8liu): should generalize this into a network device.
type Dialog struct {
	out Sender
}

// NewDialog creates a new dialog service device.
func NewDialog(out Sender) *Dialog {
	return &Dialog{
		out: out,
	}
}

// Handle handles an incoming VPC call.
func (d *Dialog) Handle(req []byte) ([]byte, int32) {
	if d.out == nil {
		return nil, 0
	}
	d.out.Send(req)
	return nil, 0
}
