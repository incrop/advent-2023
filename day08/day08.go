package day08

import (
	"advent/utils"
	"fmt"
	"regexp"
	"strings"
)

type Direction byte

const (
	LEFT Direction = iota
	RIGHT
	SIZE
)

type Node string

type Fork struct {
	from  Node
	left  Node
	right Node
}

type DesertMap struct {
	directions []Direction
	forks      map[Node]Fork
}

type State struct {
	directionIdx int
	node         Node
}

func (desertMap DesertMap) printGraphviz() {
	fmt.Println("digraph G {")
	for node, fork := range desertMap.forks {
		fmt.Printf("  \"%s\" -> \"%s\" [color=\"green\"];\n", node, fork.left)
		fmt.Printf("  \"%s\" -> \"%s\" [color=\"blue\"];\n", node, fork.right)
	}
	fmt.Println("}")
}

func (desertMap DesertMap) targetStateSteps(from Node, target NodeMatcher) (stateToSteps map[State]int) {
	stateToSteps = make(map[State]int, 0)
	step := 0
	node := from
	for {
		for idx, direction := range desertMap.directions {
			if target.matches(node) {
				state := State{
					directionIdx: idx,
					node:         node,
				}
				_, exists := stateToSteps[state]
				if exists {
					return stateToSteps
				}
				stateToSteps[state] = step
			}
			fork := desertMap.forks[node]
			switch direction {
			case LEFT:
				node = fork.left
			case RIGHT:
				node = fork.right
			}
			step++
		}
	}
}

func (desertMap DesertMap) printAllTargetStateSteps(source NodeMatcher, target NodeMatcher) {
	for sourceNode := range desertMap.forks {
		if !source.matches(sourceNode) {
			continue
		}
		fmt.Println(desertMap.targetStateSteps(sourceNode, target))
	}
}

type NodeMatcher interface {
	matches(Node) bool
}

var maxMatch int

func allMatches(matcher NodeMatcher, nodes []Node) bool {
	match := 0
	for _, node := range nodes {
		if matcher.matches(node) {
			match += 1
		}
	}
	if match >= maxMatch {
		maxMatch = match
		fmt.Println(nodes, match)
	}
	return match == len(nodes)
}

// func allMatches(matcher NodeMatcher, nodes []Node) bool {
// 	for _, node := range nodes {
// 		if !matcher.matches(node) {
// 			return false
// 		}
// 	}
// 	return true
// }

func (example Node) matches(node Node) bool {
	return example == node
}

type SuffixMatcher string

func (suffix SuffixMatcher) matches(node Node) bool {
	return strings.HasSuffix(string(node), string(suffix))
}

func (desertMap DesertMap) countDirectionsSteps(fromMatching NodeMatcher, toMatching NodeMatcher) (stepCount int) {
	curr := make([]Node, 0)
	for node := range desertMap.forks {
		if fromMatching.matches(node) {
			curr = append(curr, node)
		}
	}

	stepCount = 0
	for {
		for _, direction := range desertMap.directions {
			if allMatches(toMatching, curr) {
				return
			}
			for i, node := range curr {
				fork := desertMap.forks[node]
				switch direction {
				case LEFT:
					node = fork.left
				case RIGHT:
					node = fork.right
				}
				curr[i] = node
			}
			stepCount++
		}
	}
}

type ParsedLine struct {
	directions []Direction
	fork       *Fork
}

var directionsRe = regexp.MustCompile(`^[LR]+$`)
var forkRe = regexp.MustCompile(`^(\w+) = \((\w+), (\w+)\)$`)

func parseLine(line string) (parsed ParsedLine) {
	if directionsRe.MatchString(line) {
		parsed.directions = make([]Direction, len(line))
		for i, char := range line {
			switch char {
			case 'L':
				parsed.directions[i] = LEFT
			case 'R':
				parsed.directions[i] = RIGHT
			default:
				panic("unexpected char")
			}
		}
		return
	}
	if match := forkRe.FindStringSubmatch(line); match != nil {
		parsed.fork = &Fork{
			from:  Node(match[1]),
			left:  Node(match[2]),
			right: Node(match[3]),
		}
		return
	}
	return
}

func populateDesertMap(desertMap DesertMap, parsed ParsedLine) DesertMap {
	if parsed.directions != nil {
		desertMap.directions = parsed.directions
	}
	if parsed.fork != nil {
		if desertMap.forks == nil {
			desertMap.forks = make(map[Node]Fork)
		}
		desertMap.forks[parsed.fork.from] = *parsed.fork
	}
	return desertMap
}

func Run() int {
	desertMap := utils.ProcessInput("day08.txt", DesertMap{}, parseLine, populateDesertMap)
	// return desertMap.countDirectionsSteps(Node("AAA"), Node("ZZZ"))
	// return desertMap.countDirectionsSteps(SuffixMatcher("A"), SuffixMatcher("Z"))
	desertMap.printAllTargetStateSteps(SuffixMatcher("A"), SuffixMatcher("Z"))
	return 0
}
