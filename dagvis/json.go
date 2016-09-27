package dagvis

import (
	"encoding/json"
	"sort"
)

// N is a node in the minified DAG visualization result.
type N struct {
	X    int      `json:"x"`
	Y    int      `json:"y"`
	Ins  []string `json:"i"`
	Outs []string `json:"o"`
}

// M is a node in the minified DAG visualization result.
type M struct {
	Height int           `json:"h"`
	Width  int           `json:"w"`
	Nodes  map[string]*N `json:"n"`
}

// JSONMap returns a json'able object of a map.
func JSONMap(m *Map) interface{} {
	res := &M{
		Height: m.Height,
		Width:  m.Width,
		Nodes:  make(map[string]*N),
	}

	for name, node := range m.Nodes {
		ins := make([]string, len(node.CritIns))
		i := 0
		for in := range node.CritIns {
			ins[i] = in
			i++
		}

		outs := make([]string, len(node.CritOuts))
		i = 0
		for out := range node.CritOuts {
			outs[i] = out
			i++
		}

		sort.Strings(ins)
		sort.Strings(outs)

		n := &N{
			X:    node.X,
			Y:    node.Y,
			Ins:  ins,
			Outs: outs,
		}

		res.Nodes[name] = n
	}

	return res
}

func jsonMap(m *Map) []byte {
	res := JSONMap(m)

	ret, e := json.MarshalIndent(res, "", "    ")
	if e != nil {
		panic(e)
	}

	return ret
}
