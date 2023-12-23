package day23

import (
	"advent/utils"
	"fmt"
	"strings"
)

type Tile rune

const (
	Forest     Tile = '#'
	Path       Tile = '.'
	SlopeUp    Tile = '^'
	SlopeRight Tile = '>'
	SlopeLeft  Tile = '<'
	SlopeDown  Tile = 'v'
)

const IGNORE_SLOPES = true

type Cell struct {
	tile    Tile
	visited bool
}

type Labyrinth [][]Cell

func (labyrinth Labyrinth) String() string {
	var sb strings.Builder
	for _, row := range labyrinth {
		for _, cell := range row {
			if cell.visited {
				sb.WriteString("\033[31m")
			}
			sb.WriteRune(rune(cell.tile))
			if cell.visited {
				sb.WriteString("\033[0m")
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (labyrinth Labyrinth) findLongestPath(i, j int) (int, bool) {
	h, w := len(labyrinth), len(labyrinth[0])
	if i < 0 || i >= h || j < 0 || j >= w {
		return 0, false
	}
	cell := &labyrinth[i][j]
	if cell.visited {
		return 0, false
	}
	cell.visited = true
	defer func() {
		cell.visited = false
	}()
	switch cell.tile {
	case Forest:
		return 0, false
	case SlopeUp:
		path, found := labyrinth.findLongestPath(i-1, j)
		return path + 1, found
	case SlopeLeft:
		path, found := labyrinth.findLongestPath(i, j-1)
		return path + 1, found
	case SlopeRight:
		path, found := labyrinth.findLongestPath(i, j+1)
		return path + 1, found
	case SlopeDown:
		path, found := labyrinth.findLongestPath(i+1, j)
		return path + 1, found
	case Path:
		if i == h-1 {
			// fmt.Println(labyrinth)
			return 1, true
		}
		maxPath, maxFound := 0, false
		for _, coords := range [4]Coords{{i - 1, j}, {i, j - 1}, {i, j + 1}, {i + 1, j}} {
			path, found := labyrinth.findLongestPath(coords.i, coords.j)
			if found && path > maxPath {
				maxPath = path
				maxFound = true
			}
		}
		return maxPath + 1, maxFound
	}
	panic("unreachable")
}

type Node int

type Link struct {
	node Node
	len  int
}

type Graph struct {
	src, dst Node
	links    map[Node][]Link
}

func (graph Graph) printGraphviz() {
	fmt.Println("digraph {")
	for a, links := range graph.links {
		for _, link := range links {
			b := link.node
			fmt.Printf("%d -> %d [ label=%d ];\n", a, b, link.len)
		}
	}
	fmt.Println("}")
}

func (graph Graph) longestPath() int {
	result, found := graph.longestPathRec(
		0,
		make([]bool, len(graph.links)),
		graph.src,
	)
	if !found {
		panic("kak tak?")
	}
	return result
}

func (graph Graph) longestPathRec(acc int, visited []bool, curr Node) (int, bool) {
	if curr == graph.dst {
		return acc, true
	}
	fromIdx := int(curr)
	if visited[fromIdx] {
		return 0, false
	}
	visited[fromIdx] = true
	defer func() {
		visited[fromIdx] = false
	}()
	longestPath, longestPathFound := acc, false
	for _, link := range graph.links[curr] {
		path, found := graph.longestPathRec(acc+link.len, visited, link.node)
		if found && path > longestPath {
			longestPath, longestPathFound = path, true
		}
	}
	return longestPath, longestPathFound
}

type Coords struct {
	i, j int
}

type GraphBuilder struct {
	nodes map[Coords]Node
	links map[Node][]Link
}

func (gb *GraphBuilder) connect(a, b Node, len int) {
NEXT_PAIR:
	for _, pair := range [][2]Node{{a, b}, {b, a}} {
		node := pair[0]
		newLink := Link{pair[1], len}
		for _, existingLink := range gb.links[pair[0]] {
			if existingLink == newLink {
				continue NEXT_PAIR
			}
		}
		gb.links[node] = append(gb.links[node], newLink)
	}
}

func (gb *GraphBuilder) walk(labyrinth Labyrinth, link Link, prev, curr Coords) {
	if node, exists := gb.nodes[curr]; exists {
		gb.connect(node, link.node, link.len)
		return
	}
	i, j := curr.i, curr.j
	nexts := make([]Coords, 0, 3)
	for _, next := range []Coords{{i - 1, j}, {i, j - 1}, {i, j + 1}, {i + 1, j}} {
		if next == prev {
			continue
		}
		if next.i == len(labyrinth) {
			continue
		}
		if labyrinth[next.i][next.j].tile != Path {
			continue
		}
		nexts = append(nexts, next)
	}
	if len(nexts) == 1 {
		link.len++
		gb.walk(labyrinth, link, curr, nexts[0])
		return
	}
	node := Node(len(gb.nodes))
	gb.nodes[curr] = node
	gb.connect(node, link.node, link.len)
	for _, next := range nexts {
		gb.walk(labyrinth, Link{node, 1}, curr, next)
	}
}

func (labyrinth Labyrinth) toGraph() (graph Graph) {
	graph.src = Node(0)
	gb := &GraphBuilder{
		nodes: map[Coords]Node{{0, 1}: graph.src},
		links: make(map[Node][]Link),
	}
	gb.walk(labyrinth, Link{graph.src, 1}, Coords{0, 1}, Coords{1, 1})
	graph.dst = Node(len(gb.nodes) - 1)
	graph.links = gb.links
	return
}

func parseRow(line string) (row []Cell) {
	row = make([]Cell, len(line))
	for i, r := range line {
		tile := Tile(r)
		if IGNORE_SLOPES {
			switch tile {
			case SlopeUp, SlopeLeft, SlopeRight, SlopeDown:
				tile = Path
			}
		}
		row[i] = Cell{tile: tile}
	}
	return
}

func appendRow(labyrinth Labyrinth, row []Cell) Labyrinth {
	return append(labyrinth, row)
}

func Run() int {
	labyrinth := utils.ProcessInput("day23.txt", Labyrinth{}, parseRow, appendRow)
	fmt.Println(labyrinth)
	// path, _ := labyrinth.findLongestPath(0, 1)
	graph := labyrinth.toGraph()
	graph.printGraphviz()
	return graph.longestPath()
}
