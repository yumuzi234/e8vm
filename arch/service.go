package arch

// Request is a generic RPC request issuing from the VM.
type Request struct {
	Service uint32
	Method  uint32
	Payload []byte
}

// Service is a RPC service that handles incoming requests.
type Service interface {
	Handle(req *Request, resp []byte) (size, code uint32)
}
