package dagvis

func minCircle(nodes map[string]*MapNode) []*MapNode {
	minDist := 2 * len(nodes)
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
			cur := queue[len(queue)-1]
			count[cur]++
			if cur.Outs == nil || count[cur] > len(cur.Outs) {
				queue = queue[:len(queue)-1]
				continue
			}

			for _, next := range cur.Outs {
				if next == node {
					if len(queue) < minDist {
						minDist = len(queue)
						ret = make([]*MapNode, len(queue))
						copy(ret, queue)
					}
					break
				}
				if !flag[next] && len(queue) < minDist {
					queue = append(queue, next)
					flag[next] = true
					break
				}
			}
		}
	}

	return ret
}
