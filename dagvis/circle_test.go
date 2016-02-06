package dagvis

import (
	"testing"

	"sort"
	"strconv"
	"reflect"
)

func TestFindCircle(t *testing.T) {
	o := func(nodes map[string][]string, circExpect []string) {
		m, err := initMap(&Graph{Nodes: nodes})
		if err != nil {
			t.Fatal(err)
		}
		res := shortestCircle(m.Nodes)

		var circGot []string
		for _, node := range res {
			circGot = append(circGot, node.Name)
		}
		sort.Strings(circGot)

		if !reflect.DeepEqual(circGot, circExpect) {
			t.Errorf("min circle of %v, got %v, expect %v",
				nodes, circGot, circExpect,
			)
		}
	}

	o(map[string][]string{}, nil)
	o(map[string][]string{"a": {}}, nil)

	o(map[string][]string{
		"1": {"2"},
		"2": {"1"},
	}, []string{"1", "2"})

	o(map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {"1"},
	}, []string{"1", "2", "3"})

	o(map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {},
	}, nil)

	o(map[string][]string{
		"1": {"2"},
		"2": {"3"},
		"3": {"1", "4"},
		"4": {"5"},
		"5": {"1"},
	}, []string{"1", "2", "3"})

	nodes := map[string][]string{
		"1": {"3", "4"},
		"3": {"6", "4"},
		"4": {"3", "8"},
		"8": {"1"},
	}
	for i := 0; i < 100; i++ {
		k := strconv.Itoa(i)
		if nodes[k] != nil {
			continue
		}
		nodes[k] = []string{}
	}

	o(nodes, []string{"3", "4"})
}
