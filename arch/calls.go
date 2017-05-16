package arch

import (
	"container/list"
	"log"
	"time"

	"shanhu.io/smlvm/arch/devs"
	"shanhu.io/smlvm/net"
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
	net      net.Handler

	timedSleep bool
	sleepDur   time.Duration
}

func newCalls(p *page, mem *phyMemory, h net.Handler) *calls {
	return &calls{
		p:        &pageOffset{p, 0},
		mem:      mem,
		services: make(map[uint32]devs.Service),
		queue:    list.New(),
		net:      h,
	}
}

func (c *calls) register(id uint32, s devs.Service) {
	if id == 0 {
		panic("cannot register service 0")
	}
	c.services[id] = s
}

func (c *calls) sleep(in []byte) ([]byte, int32, *Excep) {
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
		c.sleepDur = time.Duration(Endian.Uint64(in[:8]))
		return nil, devs.ErrInternal, errSleep
	}

	// second time, timeout, waking up.
	c.timedSleep = false
	return nil, devs.ErrTimeout, nil
}

func (c *calls) system(ctrl uint8, in []byte, respSize int) (
	[]byte, int32, *Excep,
) {
	switch ctrl {
	case 1: // poll message
		if c.queue.Len() == 0 {
			return c.sleep(in)
		}

		// incoming packet queue
		front := c.queue.Front()
		p := front.Value.([]byte)
		if len(p) > respSize {
			return nil, devs.ErrSmallBuf, nil
		}
		c.queue.Remove(front)
		c.p.writeU32(callsService, 0) // a network packet
		return p, 0, nil
	case 2: // send packet out
		if c.net == nil {
			return nil, devs.ErrInvalidArg, nil
		}

		err := c.net.HandlePacket(in)
		if err != nil {
			log.Println(err)
			return nil, devs.ErrInternal, nil
		}

		return nil, 0, nil
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
	ctrl := c.p.readU8(callsControl)
	if ctrl == 0 {
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
	resp, code, exp := c.call(ctrl, service, req, respSize)
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
	return c.sleepDur, c.timedSleep
}

func (c *calls) hasPending() bool {
	return c.queue.Len() > 0
}

func (c *calls) HandlePacket(p []byte) error {
	c.queue.PushBack(p)
	return nil
}
