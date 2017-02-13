package devs

import (
	"encoding/binary"
	"time"
)

var progStartTime = time.Now()

// Clock provides time reading.
type Clock struct {
	// Now is the function pointer for reading the time.
	Now func() time.Time

	// PerfNow is a function pointer for monotonic clock.
	PerfNow func() time.Duration

	// StartTime will be used for program start time if
	// not null.
	StartTime *time.Time
}

func (c *Clock) now() time.Time {
	if c.Now == nil {
		return time.Now()
	}
	return c.Now()
}

func (c *Clock) startTime() time.Time {
	if c.StartTime != nil {
		return *c.StartTime
	}
	return progStartTime
}

func (c *Clock) perfNow() time.Duration {
	if c.PerfNow == nil {
		return time.Since(c.startTime())
	}
	return c.PerfNow()
}

func (c *Clock) read() ([]byte, int32) {
	now := c.now().UnixNano()
	ret := make([]byte, 8)
	binary.LittleEndian.PutUint64(ret, uint64(now))
	return ret, 0
}

func (c *Clock) readMono() ([]byte, int32) {
	d := c.perfNow()
	ret := make([]byte, 8)
	binary.LittleEndian.PutUint64(ret, uint64(d))
	return ret, 0
}

// Handle returns the current time as a uint64.
func (c *Clock) Handle(in []byte) ([]byte, int32) {
	if len(in) == 0 {
		return c.read()
	}

	if len(in) != 1 {
		return nil, ErrInvalidArg
	}

	switch in[0] {
	case 0:
		return c.read()
	case 1:
		return c.readMono()
	}
	return nil, ErrInvalidArg
}
