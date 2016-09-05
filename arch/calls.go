package arch

const (
	callsControl = 0x0
	callsError   = 0x1
	callsService = 0x4
	callsMethod  = 0x6

	callsRequestAddr  = 0x8
	callsRequestLen   = 0xc
	callsResponseAddr = 0x10
	callsResponseSize = 0x14 // buffer size
	callsResponseCode = 0x18
	callsResponseLen  = 0x1c

	callsSize = 0x20
)

const (
	callsErrNone = iota
	callsErrServiceNotFound
	callsErrMemoryError
)

type calls struct {
	p        *pageOffset
	mem      *phyMemory
	services map[uint32]Service
}

func newCalls(p *pageOffset, mem *phyMemory) *calls {
	return &calls{
		p:        p,
		mem:      mem,
		services: make(map[uint32]Service),
	}
}

func (c *calls) call(req *Request, resp []byte) (n, res uint32, ok bool) {
	s, found := c.services[req.Service]
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
	method := c.p.readWord(callsMethod)

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

	respLen, code, found := c.call(&Request{
		Service: service,
		Method:  method,
		Payload: req,
	}, resp)
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
