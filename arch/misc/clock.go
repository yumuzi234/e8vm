package misc

import (
	"encoding/binary"
	"time"
)

// Clock provides time reading.
type Clock struct{}

// NewClock creates a new clock
func NewClock() *Clock {
	return new(Clock)
}

// Handle returns the current time as a uint64.
func (c *Clock) Handle(_ []byte) ([]byte, int32) {
	now := time.Now().UnixNano()
	ret := make([]byte, 8)
	binary.LittleEndian.PutUint64(ret, uint64(now))
	return ret, 0
}
