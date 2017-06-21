package cluster

import (
	"shanhu.io/smlvm/arch"
)

// Cluster is a set of machine units.
type Cluster struct {
	units  []*unit
	master int
}

// Add adds a machine into the cluster.
func (c *Cluster) Add(conf *arch.Config) {

}
