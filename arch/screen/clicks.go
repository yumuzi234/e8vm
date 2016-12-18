package screen

import (
	"fmt"

	"shanhu.io/smlvm/arch/devs"
)

// ScreenClicks manages clicks on a screen.
type ScreenClicks struct {
	send devs.Sender
}

// NewClicks creates new clicks handler.
func NewClicks(s devs.Sender) *ScreenClicks {
	return &ScreenClicks{send: s}
}

// Click sends a click.
func (c *ScreenClicks) Click(line, col uint8) error {
	if line > ScreenHeight {
		return fmt.Errorf("line too big: %d", line)
	}
	if col > ScreenWidth {
		return fmt.Errorf("col too big: %d", col)
	}

	c.send.Send([]byte{byte(line), byte(col)})
	return nil
}
