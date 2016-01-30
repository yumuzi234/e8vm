package dagvis

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkFindCircle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		test()
	}
}

func test() []*MapNode {

	ret := make(map[string][]string)

	for i := 0; i < 1000; i++ {
		var edge []string
		for j := 0; j < 1000; j++ {
			if i == j {
				continue
			}
			q := rand.Int31n(500)
			if q < 1 {
				var flag bool
				for _, ele := range ret[strconv.Itoa(j)] {
					if ele == strconv.Itoa(i) {
						flag = true
					}
				}
				if !flag {
					edge = append(edge, strconv.Itoa(j))
				}
			}
		}
		ret[strconv.Itoa(i)] = edge
	}

	g := &Graph{Nodes: ret}
	nodes, _ := initMap(g)

	return shortestCircle(nodes.Nodes)
}

func TestFindCircle(t *testing.T) {
	eo := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}

	ret := make(map[string][]string)

	for i := 0; i < 1000; i++ {
		var edge []string
		if i == 1 {
			edge = append(edge, strconv.Itoa(3))
			edge = append(edge, strconv.Itoa(4))
		}
		if i == 3 {
			edge = append(edge, strconv.Itoa(6))
			edge = append(edge, strconv.Itoa(4))
		}
		if i == 4 {
			edge = append(edge, strconv.Itoa(3))
			edge = append(edge, strconv.Itoa(8))
		}
		if i == 8 {
			edge = append(edge, strconv.Itoa(1))
		}
		ret[strconv.Itoa(i)] = edge
	}

	g := &Graph{Nodes: ret}
	nodes, _ := initMap(g)

	res := shortestCircle(nodes.Nodes)

	size := len(res)
	resNode1 := res[0]
	resNode2 := res[1]

	eo(size != 2 || (resNode1.Name != "3" || resNode2.Name != "4") &&
		(resNode1.Name != "4" || resNode2.Name != "3"),
		"findSmallestCircle result error")

	for _, resNode := range res {
		fmt.Printf("%v\n", resNode.Name)
	}
}
