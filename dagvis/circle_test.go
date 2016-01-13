package dagvis

import "fmt"
import "testing"
import "strconv"
import "math/rand"

func BenchmarkTest(b *testing.B) {
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
			q := rand.Float32()
			if q < 0.002 {
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

		// if i == 1 {
		// 	edge = append(edge, strconv.Itoa(3))
		// 	edge = append(edge, strconv.Itoa(4))
		// }
		// if i == 3 {
		// 	edge = append(edge, strconv.Itoa(6))
		// 	edge = append(edge, strconv.Itoa(4))
		// }
		// if i == 4 {
		// 	edge = append(edge, strconv.Itoa(3))
		// 	edge = append(edge, strconv.Itoa(8))
		// }
		// if i == 8 {
		// 	edge = append(edge, strconv.Itoa(1))
		// }
		ret[strconv.Itoa(i)] = edge
	}

	g := &Graph{Nodes: ret}
	nodes, _ := initMap(g)

	return shortestCircle(nodes.Nodes)

}

func TestCircle(t *testing.T) {
	res := test()

	for _, resNode := range res {
		fmt.Printf("%v\n", resNode.Name)
	}
}
