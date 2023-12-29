package day25

import (
	"advent/utils"
	"fmt"
	"math/rand"
)

type Node string

func (node Node) size() int {
	return len(utils.Fields(string(node), ","))
}

type Link struct {
	a, b    Node
	enabled bool
}

type Graph struct {
	nodesToLinkIdxs map[Node][]int
	links           []Link
}

func (graph Graph) nodeSizeProduct() (product int) {
	product = 1
	for node := range graph.nodesToLinkIdxs {
		product *= node.size()
	}
	return
}

func NewGraph(nodeRows [][]Node) (graph Graph) {
	rand.Shuffle(len(nodeRows), func(i, j int) {
		nodeRows[i], nodeRows[j] = nodeRows[j], nodeRows[i]
	})
	graph.nodesToLinkIdxs = map[Node][]int{}
	for _, nodes := range nodeRows {
		a := nodes[0]
		for _, b := range nodes[1:] {
			linkIdx := len(graph.links)
			graph.links = append(graph.links, Link{a, b, true})
			graph.nodesToLinkIdxs[a] = append(graph.nodesToLinkIdxs[a], linkIdx)
			graph.nodesToLinkIdxs[b] = append(graph.nodesToLinkIdxs[b], linkIdx)
		}
	}
	return
}

func (graph Graph) neighbors(node Node) (neighbors []Node) {
	for _, linkidx := range graph.nodesToLinkIdxs[node] {
		link := graph.links[linkidx]
		if !link.enabled {
			continue
		}
		neighbor := link.a
		if neighbor == node {
			neighbor = link.b
		}
		neighbors = append(neighbors, neighbor)
	}
	return
}

func (graph Graph) visitRec(node Node, groupNum int, visited map[Node]int) {
	if visited[node] != 0 {
		return
	}
	visited[node] = groupNum
	for _, next := range graph.neighbors(node) {
		graph.visitRec(next, groupNum, visited)
	}
	return
}

func (graph Graph) sizeOfGroupsMultiplied() (count int, product int) {
	visited := map[Node]int{}
	for node := range graph.nodesToLinkIdxs {
		if visited[node] == 0 {
			count++
			graph.visitRec(node, count, visited)
		}
	}
	sizes := map[int]int{}
	for node := range graph.nodesToLinkIdxs {
		groupNum := visited[node]
		sizes[groupNum]++
	}
	product = 1
	for _, size := range sizes {
		product *= size
	}
	return
}

func (graph Graph) sizeOfGroupsAfterRemovingLinksMultiplied(n, startIdx int) (maxCount int, maxProduct int) {
	if n == 0 {
		return graph.sizeOfGroupsMultiplied()
	}
	for i := startIdx; i <= len(graph.links)-n; i++ {
		graph.links[i].enabled = false
		count, product := graph.sizeOfGroupsAfterRemovingLinksMultiplied(n-1, i+1)
		if count == 2 {
			return count, product
		}
		graph.links[i].enabled = true
	}
	return
}

// https://www.geeksforgeeks.org/introduction-and-implementation-of-kargers-algorithm-for-minimum-cut/z
func (graph Graph) findCutWithLinksNumber(n int) Graph {
	for i := 0; ; i++ {
		cut := graph.randomCut()
		linksNum := 0
		for _, link := range cut.links {
			if link.enabled {
				linksNum++
			}
		}
		if linksNum == n {
			fmt.Printf("Found after %d iterations\n", i)
			return cut
		}
	}
}

func (graph Graph) randomCut() (cut Graph) {
	cut.nodesToLinkIdxs = make(map[Node][]int)
	for node, linkIdx := range graph.nodesToLinkIdxs {
		cut.nodesToLinkIdxs[node] = linkIdx
	}
	cut.links = make([]Link, len(graph.links))
	copy(cut.links, graph.links)
	enabledLinks := len(graph.links)
	for len(cut.nodesToLinkIdxs) > 2 {
		var a, b Node
		n := rand.Intn(enabledLinks)
		for _, link := range cut.links {
			if !link.enabled {
				continue
			}
			if n == 0 {
				a, b = link.a, link.b
				break
			}
			n--
		}
		abNode := a + "," + b
		for _, a_or_b := range []Node{a, b} {
			for _, linkIdxs := range cut.nodesToLinkIdxs[a_or_b] {
				link := &cut.links[linkIdxs]
				if link.a == a_or_b {
					link.a = abNode
				}
				if link.b == a_or_b {
					link.b = abNode
				}
				if link.a == link.b {
					link.enabled = false
					enabledLinks--
				}
			}
		}
		abLinks := make([]int, 0)
		for idx, link := range cut.links {
			if link.enabled && (link.a == abNode || link.b == abNode) {
				abLinks = append(abLinks, idx)
			}
		}
		delete(cut.nodesToLinkIdxs, a)
		delete(cut.nodesToLinkIdxs, b)
		cut.nodesToLinkIdxs[abNode] = abLinks
	}
	return
}

func parseNodes(line string) (nodes []Node) {
	fields := utils.Fields(line, ": ")
	nodes = make([]Node, len(fields))
	for i, field := range fields {
		nodes[i] = Node(field)
	}
	return
}

func aggregate(nodeRows [][]Node, nodes []Node) [][]Node {
	return append(nodeRows, nodes)
}

func Run() int {
	nodeRows := utils.ProcessInput("day25.txt", nil, parseNodes, aggregate)
	graph := NewGraph(nodeRows)
	return graph.findCutWithLinksNumber(3).nodeSizeProduct()
}
