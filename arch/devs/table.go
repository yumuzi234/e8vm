package devs

// Table is a virtual card table device.
type Table struct {
	in Sender
}

// NewTable creates a new virtual card table device.
func NewTable(in Sender) *Table {
	return &Table{
		in: in,
	}
}

func whatCode(what string) uint8 {
	switch what {
	case "card":
		return 1
	case "button":
		return 2
	case "div":
		return 3
	case "box":
		return 4
	}
	return 0
}

// Click sends in a click on the table.
func (t *Table) Click(what string, pos uint8) {
	t.in.Send([]byte{whatCode(what), pos})
}
