package devs

// Service is a RPC service that handles incoming requests.
type Service interface {
	Handle(req []byte) ([]byte, int32)
}
