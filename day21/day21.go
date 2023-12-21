package day21

import (
	"advent/utils"
	"fmt"
	"strings"
)

type Tile rune

type Position struct {
	i, j int
}

type State struct {
	tiles     [][]Tile
	positions map[Position]bool
}

func (state State) walk(steps int) State {
	h, w := len(state.tiles), len(state.tiles[0])
	for n := 0; n < steps; n++ {
		next := make(map[Position]bool)
		for pos := range state.positions {
			i, j := pos.i, pos.j
			if i > 0 && state.tiles[i-1][j] == '.' {
				next[Position{i - 1, j}] = true
			}
			if i < h-1 && state.tiles[i+1][j] == '.' {
				next[Position{i + 1, j}] = true
			}
			if j > 0 && state.tiles[i][j-1] == '.' {
				next[Position{i, j - 1}] = true
			}
			if j < w-1 && state.tiles[i][j+1] == '.' {
				next[Position{i, j + 1}] = true
			}
		}
		state.positions = next
	}
	return state
}

func (state State) posCount() int {
	return len(state.positions)
}

func (state State) reset(i, j int) State {
	return State{state.tiles, map[Position]bool{{i, j}: true}}
}

func (state State) String() string {
	var sb strings.Builder
	for i, row := range state.tiles {
		for j, tile := range row {
			if state.positions[Position{i, j}] {
				sb.WriteRune('O')
			} else {
				sb.WriteRune(rune(tile))
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (state State) print() State {
	fmt.Println(state)
	return state
}

func parseRow(line string) (row []Tile) {
	row = make([]Tile, len(line))
	for i, r := range line {
		row[i] = Tile(r)
	}
	return
}

func aggregate(state State, row []Tile) State {
	if state.positions == nil {
		for j, t := range row {
			if t == Tile('S') {
				row[j] = Tile('.')
				i := len(state.tiles)
				state.positions = map[Position]bool{{i, j}: true}
				break
			}
		}
	}
	state.tiles = append(state.tiles, row)
	return state
}

func (state State) infiniteWalkSteps(R int) (total int) {
	size := len(state.tiles)
	for i, row := range state.tiles {
		expandedRow := make([]Tile, size*(R*2+1))
		for n := 0; n < R*2+1; n++ {
			copy(expandedRow[n*size:], row)
		}
		state.tiles[i] = expandedRow
	}
	expandedTiles := make([][]Tile, size*(R*2+1))
	for i, row := range state.tiles {
		for n := 0; n < R*2+1; n++ {
			expandedTiles[i+n*size] = make([]Tile, len(row))
			copy(expandedTiles[i+n*size], row)
		}
	}

	state.tiles = expandedTiles
	center := R*size + size/2
	return state.reset(center, center).walk(center).posCount()
}

func (state State) infiniteWalkPosCountOptimized(steps int) (total int) {
	size := len(state.tiles)
	center := size / 2
	R := (steps - center) / size
	{
		fullCount := [2]int{1, 0}
		idx := 1
		incr := 4
		for i := 1; i < R; i++ {
			fullCount[idx] += incr
			incr += 4
			idx = 1 - idx
		}
		state := state.reset(center, center).walk(size + (R % 2))
		total += fullCount[0] * state.posCount()
		total += fullCount[1] * state.walk(1).posCount()
	}
	{
		state := state.reset(0, 0).walk(center - 1)
		total += R * state.posCount()
		total += (R - 1) * state.walk(size).posCount()
	}
	total += state.reset(0, center).walk(size - 1).posCount()
	{
		state := state.reset(0, size-1).walk(center - 1)
		total += R * state.posCount()
		total += (R - 1) * state.walk(size).posCount()
	}
	total += state.reset(center, size-1).walk(size - 1).posCount()
	{
		state := state.reset(size-1, size-1).walk(center - 1)
		total += R * state.posCount()
		total += (R - 1) * state.walk(size).posCount()
	}
	total += state.reset(size-1, center).walk(size - 1).posCount()
	{
		state := state.reset(size-1, 0).walk(center - 1)
		total += R * state.posCount()
		total += (R - 1) * state.walk(size).posCount()
	}
	total += state.reset(center, 0).walk(size - 1).posCount()
	return
}

func Run() int {
	state := utils.ProcessInput("day21.txt", State{}, parseRow, aggregate)
	return state.infiniteWalkPosCountOptimized(26501365)
}
