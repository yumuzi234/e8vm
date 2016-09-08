package arch

import (
	"fmt"

	"e8vm.io/e8vm/arch/vpc"
)

type click struct {
	line uint8
	col  uint8
}

type clicks struct {
	send vpc.Sender
}

func newClicks(s vpc.Sender) *clicks {
	return &clicks{
		send: s,
	}
}

func (c *clicks) addClick(line, col uint8) error {
	if line > screenHeight {
		return fmt.Errorf("line too big: %d", line)
	}
	if col > screenWidth {
		return fmt.Errorf("col too big: %d", col)
	}

	c.send.Send([]byte{byte(line), byte(col)})
	return nil
}
