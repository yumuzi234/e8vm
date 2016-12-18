package devs

import (
	crand "crypto/rand"
	"io"
	"log"
	"math/rand"
	"time"
)

// Rand provides a random number generator.
type Rand struct {
	r io.Reader
}

// NewRand creates a new random number generator.
func NewRand(seed int64) *Rand {
	r := rand.New(rand.NewSource(seed))
	return &Rand{r: r}
}

// NewTimeRand creates a new random number generator with the current time as
// the seed.
func NewTimeRand() *Rand {
	s := time.Now().UnixNano()
	return NewRand(s)
}

// NewCryptoRand creates a new random number generator that uses the
// crypto/rand.Reader as the source.
func NewCryptoRand() *Rand {
	return &Rand{r: crand.Reader}
}

// Handle handles incoming request to generate a random number.
func (r *Rand) Handle(_ []byte) ([]byte, int32) {
	ret := make([]byte, 4)
	_, err := r.r.Read(ret)
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return ret, 0
}
