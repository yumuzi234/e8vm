package dagvis

func shortestCircle(nodes map[string]*MapNode) []*MapNode {
	dists := make(map[*MapNode]map[*MapNode]int)
	nexts := make(map[*MapNode]map[*MapNode]*MapNode)

	for _, node := range nodes {
		dists[node] = make(map[*MapNode]int)
		nexts[node] = make(map[*MapNode]*MapNode)
	}

	for _, from := range nodes {
		for _, to := range from.Outs {
			dists[from][to] = 1
			nexts[from][to] = to
		}
	}

	dist := func(from, to *MapNode) (d int, inf bool) {
		d, ok := dists[from][to]
		if !ok {
			return 0, true
		}
		return d, false
	}

	var shortestNode *MapNode
	var shortestDist int
	shortestDist = 2 * len(nodes)

	for _, via := range nodes {
		for _, from := range nodes {
			if from == via {
				continue
			}

			for _, to := range nodes {
				if to == via {
					continue
				}

				d1, d1Inf := dist(from, via)
				d2, d2Inf := dist(via, to)
				// if one edge is larger than min circle dist, then continue
				if d1Inf || d2Inf ||
					d1 >= shortestDist ||
					d2 >= shortestDist {
					continue
				}

				d, inf := dist(from, to)
				if inf || d1+d2 < d {
					dists[from][to] = d1 + d2
					nexts[from][to] = nexts[from][via]
				}

				if from == to {
					// check if there are two vertices form a circle, then
					// return
					if d1 == 1 && d2 == 1 {
						var ret []*MapNode
						ret = append(ret, from)
						ret = append(ret, via)
						return ret
					}
					// update current min circle distance
					if dists[from][to] < shortestDist {
						shortestDist = dists[from][to]
						shortestNode = from
					}

				}
			}
		}
	}

	if shortestNode == nil { // no circle
		return nil
	}

	var ret []*MapNode
	node := shortestNode
	for {
		ret = append(ret, node)
		if len(ret) > len(nodes) {
			panic("too big")
		}
		node = nexts[node][shortestNode]
		if node == shortestNode {
			break
		}
	}
	return ret
}
