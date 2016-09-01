package builds

import (
	"fmt"
)

func logln(c *context, s string) {
	if c.LogLine == nil {
		fmt.Println(s)
	} else {
		c.LogLine(s)
	}
}
