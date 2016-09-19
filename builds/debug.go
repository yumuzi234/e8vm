package builds

import (
	"fmt"
	"math"

	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/image"
)

func debugSection(tab *debug.Table) (*image.Section, error) {
	bs := tab.Marshal()
	if len(bs) > math.MaxInt32-1 {
		return nil, fmt.Errorf("debug section too large")
	}

	return &image.Section{
		Header: &image.Header{
			Type: image.Debug,
			Size: uint32(len(bs)),
		},
		Bytes: bs,
	}, nil
}
