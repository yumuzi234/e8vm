package screen

import (
	"fmt"

	"shanhu.io/smlvm/arch/vpc"
)

// Clicks manages clicks on a screen.
type Clicks struct {
	send vpc.Sender
}

// NewClicks creates new clicks handler.
func NewClicks(s vpc.Sender) *Clicks {
	return &Clicks{send: s}
}

// Click sends a click.
func (c *Clicks) Click(line, col uint8) error {
	if line > Height {
		return fmt.Errorf("line too big: %d", line)
	}
	if col > Width {
		return fmt.Errorf("col too big: %d", col)
	}

	c.send.Send([]byte{byte(line), byte(col)})
	return nil
}
