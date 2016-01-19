package build8

import (
	"fmt"
	"math"

	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/e8"
)

func debugSection(tab *debug8.Table) (*e8.Section, error) {
	bs := tab.Marshal()
	if len(bs) > math.MaxInt32-1 {
		return nil, fmt.Errorf("debug section too large")
	}

	return &e8.Section{
		Header: &e8.Header{
			Type: e8.Debug,
			Size: uint32(len(bs)),
		},
		Bytes: bs,
	}, nil
}
