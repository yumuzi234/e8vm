package main

import "fmt"
import "testing"

func BenchmarkTest(b *testing.B) {

	// run the test function b.N times
        for n := 0; n < b.N; n++ {
                test()
        }
 }

func test() []*MapNode {

	nodes := make (map[string]*MapNode)
	
	var a,b,c,d,e,f,g *MapNode
	a = newMapNode("a")
	b = newMapNode("b")
	c = newMapNode("c")
	d = newMapNode("d")
	e = newMapNode("e")
	f = newMapNode("f")
	g = newMapNode("g")

	// set edges
	a.Outs["ac"] = c
	a.Outs["ae"] = e
	b.Outs["bc"] = c
	b.Outs["bd"] = d
	b.Outs["be"] = e
	c.Outs["cd"] = d
	d.Outs["da"] = a
	d.Outs["de"] = e
	//e.Outs["ea"] = a

	var nodesArray []*MapNode

	nodesArray = append(nodesArray, a, b, c, d, e, f, g)
	
	for _, node := range nodesArray {
		nodes[node.Name] = node
	} 

	return shortestCircle(nodes)

 }

 func TestCircle(t *testing.T) {

	var res []*MapNode
	res = test()

	 for _, resNode := range res {	
	 	fmt.Printf("%v\n", resNode.Name)
	 } 

 }
