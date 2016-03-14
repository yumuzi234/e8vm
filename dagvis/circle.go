package dagvis

type searchNode struct {
	start  *MapNode
	this   *MapNode
	last   *searchNode
	length int
}

func traceCircle(trace []*searchNode, snode *searchNode) []*MapNode {
	n := snode.length
	ret := make([]*MapNode, n)
	for i := 0; i < n; i++ {
		ret[n-1-i] = snode.this
		snode = snode.last
	}

	if snode != nil {
		panic("bug")
	}
	return ret
}

func minCircle(nodes map[string]*MapNode) []*MapNode {
	var trace []*searchNode
	visited := make(map[string]map[string]bool)
	for _, node := range nodes {
		m := make(map[string]bool)
		m[node.Name] = true
		visited[node.Name] = m
	}

	for _, node := range nodes {
		trace = append(trace, &searchNode{
			start:  node,
			this:   node,
			last:   nil,
			length: 1,
		})
	}

	pt := 0
	for pt < len(trace) {
		snode := trace[pt]
		start := snode.start
		vmap := visited[start.Name]
		for name, out := range snode.this.Outs {
			if name == start.Name {
				return traceCircle(trace, snode)
			}

			if vmap[name] {
				// visited before from this start
				continue
			}
			if name > start.Name {
				// a node larger than the start; skip it
				continue
			}

			trace = append(trace, &searchNode{
				start:  start,
				this:   out,
				last:   snode,
				length: snode.length + 1,
			})
		}

		pt++ // next one
	}

	return nil // no circle found
}
