package day17

import (
	"advent/utils"
	"container/heap"
	"fmt"
)

type Direction byte

const (
	Right Direction = iota
	Down
	Left
	Up
)

const DIR_COUNT = 4

func (dir Direction) toRune() rune {
	switch dir {
	case Right:
		return '→'
	case Down:
		return '↓'
	case Left:
		return '←'
	case Up:
		return '↑'
	}
	return '?'
}

func (dir Direction) turn(d int) Direction {
	n := int(dir) + d
	n = ((n % DIR_COUNT) + DIR_COUNT) % DIR_COUNT
	return Direction(n)
}

func (dir Direction) offset(i, j int) (int, int) {
	switch dir {
	case Right:
		return i, j + 1
	case Down:
		return i + 1, j
	case Left:
		return i, j - 1
	case Up:
		return i - 1, j
	}
	return i, j
}

const MIN_DIR_STEPS = 4
const MAX_DIR_STEPS = 10

type Step struct {
	prev      *Step
	from      *Tile
	to        *Tile
	dir       Direction
	dirCount  int
	totalLoss int
}

func (step *Step) String() string {
	if step.from == nil {
		return fmt.Sprintf("%c (%d): %d", step.dir.toRune(), step.dirCount, step.totalLoss)
	}
	return fmt.Sprintf("[%d %d] %c (%d): %d", step.from.i, step.from.j, step.dir.toRune(), step.dirCount, step.totalLoss)
}

func (step *Step) printPath() {
	if step == nil {
		return
	}
	step.prev.printPath()
	fmt.Println(step)
}

type Tile struct {
	i, j int
	loss int
}

type Field [][]Tile

func (field Field) nextSteps(step *Step) (results []*Step) {
	results = make([]*Step, 0, 3)
	for d := -1; d <= 1; d++ {
		next := &Step{prev: step, from: step.to}
		if d == 0 {
			if step.dirCount >= MAX_DIR_STEPS {
				continue
			}
			next.dirCount = step.dirCount + 1
		} else {
			if step.dirCount < MIN_DIR_STEPS {
				continue
			}
			next.dirCount = 1
		}

		next.dir = step.dir.turn(d)
		i, j := next.dir.offset(next.from.i, next.from.j)
		if i < 0 || i >= len(field) || j < 0 || j >= len(field[0]) {
			continue
		}
		next.to = &field[i][j]
		next.totalLoss = step.totalLoss + next.to.loss
		results = append(results, next)
	}
	return
}

func (field Field) printPath(path *Step) {
	runes := make([][]rune, len(field))
	for i, row := range field {
		runes[i] = make([]rune, len(field[0]))
		for j, tile := range row {
			runes[i][j] = '0' + rune(tile.loss)
		}
	}
	for path != nil {
		runes[path.to.i][path.to.j] = path.dir.toRune()
		path = path.prev
	}
	for _, row := range runes {
		for _, r := range row {
			fmt.Printf("%c", r)
		}
		fmt.Println()
	}
	fmt.Println()
}

type PathKey struct {
	i, j     int
	dir      Direction
	dirCount int
}

func (step *Step) toKey() PathKey {
	return PathKey{
		i:        step.to.i,
		j:        step.to.j,
		dir:      step.dir,
		dirCount: step.dirCount,
	}
}

func (field Field) calculateBestPath() *Step {
	minPaths := make(map[PathKey]*Step)
	steps := &StepHeap{
		&Step{to: &field[0][0], dir: Right, dirCount: MAX_DIR_STEPS},
		&Step{to: &field[0][0], dir: Down, dirCount: MAX_DIR_STEPS},
	}
	heap.Init(steps)
	for {
		step := heap.Pop(steps).(*Step)
		if step.to.i == len(field)-1 && step.to.j == len(field[0])-1 {
			if step.dirCount >= MIN_DIR_STEPS {
				return step
			}
		}
		if minPath := minPaths[step.toKey()]; minPath != nil && minPath.totalLoss < step.totalLoss {
			continue
		}
		for _, next := range field.nextSteps(step) {
			if minPath := minPaths[next.toKey()]; minPath != nil && minPath.totalLoss <= next.totalLoss {
				continue
			}
			minPaths[next.toKey()] = next
			heap.Push(steps, next)
		}
	}
}

type StepHeap []*Step

func (h StepHeap) Len() int           { return len(h) }
func (h StepHeap) Less(i, j int) bool { return h[i].totalLoss < h[j].totalLoss }
func (h StepHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *StepHeap) Push(x interface{}) {
	*h = append(*h, x.(*Step))
}
func (h *StepHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func parseRow(line string, i int) (row []Tile) {
	row = make([]Tile, len(line))
	for j, r := range line {
		row[j] = Tile{
			i:    i,
			j:    j,
			loss: int(r - '0'),
		}
	}
	return
}

func appendRows(field Field, row []Tile) Field {
	return append(field, row)
}

func Run() int {
	field := utils.ProcessInputWithLineNumbers("day17.txt", Field{}, parseRow, appendRows)
	field.printPath(nil)
	path := field.calculateBestPath()
	field.printPath(path)
	return path.totalLoss
}
