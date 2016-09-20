package vpc

// VPC error codes.
const (
	ErrNotFound = 1 + iota
	ErrInvalidArg
	ErrMemory
	ErrSmallBuf
	ErrInternal
	ErrTimeout
)
