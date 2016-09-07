package arch

import (
	"container/list"

	"e8vm.io/e8vm/arch/vpc"
	"e8vm.io/e8vm/coder"
)

const (
	callsControl = 0x0
	callsError   = 0x1
	callsService = 0x4

	callsRequestAddr  = 0x8
	callsRequestLen   = 0xc
	callsResponseAddr = 0x10
	callsResponseSize = 0x14 // buffer size
	callsResponseCode = 0x18
	callsResponseLen  = 0x1c

	callsSize = 0x20
)

// error codes
const (
	callsErrNone = iota
	callsErrServiceNotFound
	callsErrMemoryError
)

// result codes
const (
	callsResOK = iota
	callsResInvalidRequest
)

type calls struct {
	p        *pageOffset
	mem      *phyMemory
	services map[uint32]vpc.Service
	enabled  map[uint32]bool
	queue    *list.List
}

func newCalls(p *page, mem *phyMemory) *calls {
	return &calls{
		p:        &pageOffset{p, 0},
		mem:      mem,
		services: make(map[uint32]vpc.Service),
		queue:    list.New(),
	}
}

func (c *calls) sender(id uint32) vpc.Sender {
	return &callsSender{service: id, queue: c.queue}
}

func (c *calls) register(id uint32, s vpc.Service) {
	if id == 0 {
		panic("cannot register service 0")
	}

	c.services[id] = s
}

func (c *calls) callControl(req, resp []byte) (n, res uint32) {
	dec := coder.NewDecoder(req)
	cmd := dec.U32()
	if dec.Err != nil {
		return 0, callsResInvalidRequest
	}

	// TODO
	switch cmd {
	case 0: // poll message

	case 1: // list services

	case 2: // enable service message

	case 3: // disable service message

	}

	return 0, 0
}

func (c *calls) call(service uint32, req, resp []byte) (
	n, res uint32, ok bool,
) {
	if service == 0 {
		n, res = c.callControl(req, resp)
		return n, res, false
	}

	s, found := c.services[service]
	if !found {
		return 0, 0, false
	}

	n, res = s.Handle(req, resp)
	return n, res, true
}

func (c *calls) Tick() {
	control := c.p.readByte(callsControl)
	if control == 0 {
		return
	}

	service := c.p.readWord(callsService)

	reqAddr := c.p.readWord(callsRequestAddr)
	reqLen := c.p.readWord(callsRequestLen)
	respAddr := c.p.readWord(callsResponseAddr)
	respSize := c.p.readWord(callsResponseSize)

	var req, resp []byte
	if reqLen > 0 {
		req = make([]byte, reqLen)
	}
	if respSize > 0 {
		resp = make([]byte, respSize)
	}

	for i := range req {
		var exp *Excep
		req[i], exp = c.mem.ReadByte(reqAddr + uint32(i))
		if exp != nil {
			c.p.writeByte(callsError, callsErrMemoryError)
			return
		}
	}

	respLen, code, found := c.call(service, req, resp)
	if !found {
		c.p.writeByte(callsError, callsErrServiceNotFound)
		return
	}

	if respLen > respSize {
		respLen = respSize
	}
	if resp != nil {
		resp = resp[:respLen]
	}
	for i := range resp {
		exp := c.mem.WriteByte(respAddr+uint32(i), resp[i])
		if exp != nil {
			c.p.WriteByte(callsError, callsErrMemoryError)
			return
		}
	}

	c.p.writeWord(callsResponseCode, code)
	c.p.writeWord(callsResponseLen, respLen)
	c.p.writeByte(callsControl, 0)
}
