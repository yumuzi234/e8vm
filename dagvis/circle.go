package dagvis

func minCircle2(nodes map[string]*MapNode) []*MapNode {
	panic("todo")
}

func minCircle(nodes map[string]*MapNode) []*MapNode {
	minDist := 2 * len(nodes)
	var ret []*MapNode

	for _, node := range nodes {
		if node.Outs == nil {
			continue
		}
		flag := make(map[*MapNode]bool)
		count := make(map[*MapNode]int)

		var stack []*MapNode
		stack = append(stack, node)
		flag[node] = true

		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			count[cur]++

			// why count[cur] can be compared with len(cur.Outs) ?
			if cur.Outs == nil || count[cur] > len(cur.Outs) {
				stack = stack[:len(stack)-1]
				continue
			}

			for _, next := range cur.Outs {
				if next == node {
					if len(stack) < minDist {
						minDist = len(stack)
						ret = make([]*MapNode, len(stack))
						copy(ret, stack)
					}
					break
				}
				if !flag[next] && len(stack) < minDist {
					stack = append(stack, next)
					flag[next] = true
					break
				}
			}
		}
	}

	return ret
}
