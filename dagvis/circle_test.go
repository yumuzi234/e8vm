package dagvis

import (
	"sort"
	"strconv"
	"testing"
)

func TestFindCircle(t *testing.T) {
	eo := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}
	helper := func(ret map[string][]string) (MapNodeSlice, int) {
		g := &Graph{Nodes: ret}
		nodes, _ := initMap(g)
		res := shortestCircle(nodes.Nodes)
		size := len(res)
		sort.Sort(res)
		return res, size
	}

	ret := make(map[string][]string)
	res, size := helper(ret)
	eo(res != nil, "findSmallestCircle result error")

	var edge []string
	ret[strconv.Itoa(1)] = edge
	res, size = helper(ret)
	eo(res != nil, "findSmallestCircle result error")

	for i := 1; i <= 2; i++ {
		var edge []string
		if i == 1 {
			edge = append(edge, strconv.Itoa(2))
		}
		if i == 2 {
			edge = append(edge, strconv.Itoa(1))
		}
		ret[strconv.Itoa(i)] = edge
	}
	res, size = helper(ret)
	eo(size != 2 || res[0].Name != "1" || res[1].Name != "2",
		"findSmallestCircle result error")

	for i := 1; i <= 3; i++ {
		var edge []string
		if i == 1 {
			edge = append(edge, strconv.Itoa(2))
		}
		if i == 2 {
			edge = append(edge, strconv.Itoa(3))
		}
		if i == 3 {
			edge = append(edge, strconv.Itoa(1))
		}
		ret[strconv.Itoa(i)] = edge
	}
	res, size = helper(ret)
	eo(size != 3 || res[0].Name != "1" || res[1].Name != "2" ||
		res[2].Name != "3", "findSmallestCircle result error1")

	for i := 1; i <= 3; i++ {
		var edge []string
		if i == 1 {
			edge = append(edge, strconv.Itoa(2))
		}
		if i == 2 {
			edge = append(edge, strconv.Itoa(3))
		}
		ret[strconv.Itoa(i)] = edge
	}
	res, size = helper(ret)
	eo(res != nil, "findSmallestCircle result error")

	for i := 1; i <= 5; i++ {
		var edge []string
		if i == 1 {
			edge = append(edge, strconv.Itoa(2))
		}
		if i == 2 {
			edge = append(edge, strconv.Itoa(3))
		}
		if i == 3 {
			edge = append(edge, strconv.Itoa(1))
			edge = append(edge, strconv.Itoa(4))
		}
		if i == 4 {
			edge = append(edge, strconv.Itoa(5))
		}
		if i == 5 {
			edge = append(edge, strconv.Itoa(1))
		}
		ret[strconv.Itoa(i)] = edge
	}
	res, size = helper(ret)
	eo(size != 3 || res[0].Name != "1" || res[1].Name != "2" ||
		res[2].Name != "3", "findSmallestCircle result error1")

	for i := 0; i < 100; i++ {
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
	res, size = helper(ret)
	eo(size != 2 || res[0].Name != "3" || res[1].Name != "4",
		"findSmallestCircle result error")
}
