package vpc

// Service is a RPC service that handles incoming requests.
type Service interface {
	Handle(req, resp []byte) (n, res uint32)
}
