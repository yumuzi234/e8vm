package dagvis

func shortestCircle(nodes map[string]*MapNode) MapNodeSlice {

	shortestDist := 2 * len(nodes)
	var ret []*MapNode

	for _, node := range nodes {
		if node.Outs == nil {
			continue
		}
		flag := make(map[*MapNode]bool)
		count := make(map[*MapNode]int)
		var queue []*MapNode
		queue = append(queue, node)
		flag[node] = true
		for len(queue) > 0 {
			curr := queue[len(queue)-1]
			count[curr]++
			if curr.Outs == nil || count[curr] > len(curr.Outs) {
				queue = queue[:len(queue)-1]
				continue
			}
			for _, next := range curr.Outs {
				if next == node {
					if len(queue) < shortestDist {
						shortestDist = len(queue)
						ret = make([]*MapNode, len(queue))
						copy(ret, queue)
					}
					break
				}
				if !flag[next] && len(queue) < shortestDist {
					queue = append(queue, next)
					flag[next] = true
					break
				}
			}
		}
	}

	return ret
}
