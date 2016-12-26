package devs

// Table is a virtual card table device.
type Table struct {
	out TableRender
	in  Sender
}

// NewTable creates a new virtual card table device.
func NewTable(out TableRender, in Sender) *Table {
	return &Table{
		out: out,
		in:  in,
	}
}

// Handle handles an incoming VPC.
func (t *Table) Handle(req []byte) ([]byte, int32) {
	if t.out == nil {
		return nil, 0
	}
	t.out.Act(&TableAction{req})
	return nil, 0
}

func whatCode(what string) uint8 {
	switch what {
	case "card":
		return 1
	case "button":
		return 2
	case "div":
		return 3
	}
	return 0
}

// Click sends in a click on the table.
func (t *Table) Click(what string, pos uint8) error {
	t.in.Send([]byte{whatCode(what), pos})
	return nil
}
