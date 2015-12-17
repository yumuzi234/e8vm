package dagvis

import "fmt"
import "testing"

func BenchmarkTest(b *testing.B) {

	// run the test function b.N times
	for n := 0; n < b.N; n++ {
		test()
	}
}

func test() []*MapNode {

	g := &Graph{
		Nodes: map[string][]string{
			"a": {"c", "e"},
			"b": {"c", "d", "e"},
			"c": {"d"},
			"d": {"a", "e"},
			"e": {"f"},
			"f": {},
		},
	}

	nodes, _ := initMap(g)

	return shortestCircle(nodes.Nodes)

}

func TestCircle(t *testing.T) {

	res := test()

	for _, resNode := range res {
		fmt.Printf("%v\n", resNode.Name)
	}

}
