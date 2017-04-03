package arch

import (
	"container/list"
	"time"

	"shanhu.io/smlvm/arch/devs"
)

const (
	callsControl = 0x0
	callsService = 0x4

	callsRequestAddr  = 0x8
	callsRequestLen   = 0xc
	callsResponseAddr = 0x10
	callsResponseSize = 0x14 // buffer size
	callsResponseCode = 0x18
	callsResponseLen  = 0x1c

	callsSize = 0x20
)

type calls struct {
	p        *pageOffset
	mem      *phyMemory
	services map[uint32]devs.Service
	enabled  map[uint32]bool
	queue    *list.List

	timedSleep bool
	sleep      time.Duration
}

func newCalls(p *page, mem *phyMemory) *calls {
	return &calls{
		p:        &pageOffset{p, 0},
		mem:      mem,
		services: make(map[uint32]devs.Service),
		queue:    list.New(),
	}
}

func (c *calls) sender(id uint32) devs.Sender {
	return &callsSender{service: id, queue: c.queue}
}

func (c *calls) register(id uint32, s devs.Service) {
	if id == 0 {
		panic("cannot register service 0")
	}
	c.services[id] = s
}

func (c *calls) system(ctrl uint8, in []byte, respSize int) (
	[]byte, int32, *Excep,
) {
	switch ctrl {
	case 1: // poll message
		if c.queue.Len() == 0 {
			if len(in) == 0 {
				c.timedSleep = false
				return nil, devs.ErrInternal, errSleep // we will execute again
			}
			if len(in) != 8 {
				return nil, devs.ErrInvalidArg, nil
			}

			// first time executing sleep.
			if !c.timedSleep {
				c.timedSleep = true
				c.sleep = time.Duration(Endian.Uint64(in[:8]))
				return nil, devs.ErrInternal, errSleep
			}

			// second time, timeout, waking up.
			c.timedSleep = false
			return nil, devs.ErrTimeout, nil
		}

		front := c.queue.Front()
		m := front.Value.(*callsMessage)
		if len(m.p) > respSize {
			return nil, devs.ErrSmallBuf, nil
		}

		c.queue.Remove(front)
		c.p.writeU32(callsService, m.service) // overwrite the service
		return m.p, 0, nil
	}

	return nil, devs.ErrInvalidArg, nil
}

func (c *calls) call(ctrl uint8, s uint32, req []byte, respSize int) (
	[]byte, int32, *Excep,
) {
	if s == 0 {
		return c.system(ctrl, req, respSize)
	}

	service, found := c.services[s]
	if !found {
		return nil, devs.ErrNotFound, nil
	}
	resp, ret := service.Handle(req)
	return resp, ret, nil
}

func (c *calls) respondCode(code int32) {
	c.p.writeU32(callsResponseCode, uint32(code))
}

func (c *calls) respSize() int {
	ret := c.p.readU32(callsResponseSize)
	if ret > devs.MaxLen {
		ret = devs.MaxLen
	}
	return int(ret)
}

func (c *calls) invoke() *Excep {
	control := c.p.readU8(callsControl)
	if control == 0 {
		return nil
	}

	service := c.p.readU32(callsService)
	reqAddr := c.p.readU32(callsRequestAddr)
	reqLen := c.p.readU32(callsRequestLen)

	var req []byte
	if reqLen > 0 {
		req = make([]byte, reqLen)
	}

	for i := range req {
		var exp *Excep
		req[i], exp = c.mem.ReadU8(reqAddr + uint32(i))
		if exp != nil {
			return exp
		}
	}

	respSize := c.respSize()
	resp, code, exp := c.call(control, service, req, respSize)
	if exp != nil {
		return exp
	}
	if code != 0 {
		c.respondCode(code)
		return nil
	}

	respAddr := c.p.readU32(callsResponseAddr)
	respLen := len(resp)
	if respLen > devs.MaxLen {
		c.respondCode(devs.ErrInternal)
		return nil
	}

	// we will write the response length anyways
	c.p.writeU32(callsResponseLen, uint32(respLen))
	if respLen > respSize {
		c.respondCode(devs.ErrSmallBuf)
		return nil
	}

	if resp != nil {
		for i := range resp {
			if exp := c.mem.WriteU8(respAddr+uint32(i), resp[i]); exp != nil {
				return exp
			}
		}
	}

	c.respondCode(0)
	c.p.writeU8(callsControl, 0)
	return nil
}

func (c *calls) sleepTime() (time.Duration, bool) {
	return c.sleep, c.timedSleep
}

func (c *calls) queueLen() int { return c.queue.Len() }
