package day16

import (
	"advent/utils"
	"fmt"
	"strings"
)

type Tile rune

const (
	Empty       Tile = '.'
	MirrirSlash Tile = '/'
	MirrorBack  Tile = '\\'
	SplitHor    Tile = '-'
	SplitVer    Tile = '|'
)

type Field [][]Tile

type Direction byte

const (
	Up Direction = iota
	Left
	Right
	Down
)

type Energized [4]bool

type State [][]Energized

type Game struct {
	field       Field
	isEnergized State
}

func (game Game) h() int {
	return len(game.field)
}

func (game Game) w() int {
	return len(game.field[0])
}

func (game Game) String() string {
	var sb strings.Builder
	for i, row := range game.field {
		for j, tile := range row {
			anyEnergized := false
			for _, dirEnergized := range game.isEnergized[i][j] {
				anyEnergized = anyEnergized || dirEnergized
			}
			if anyEnergized {
				sb.WriteString("\033[32m")
				sb.WriteRune(rune(tile))
				sb.WriteString("\033[0m")
			} else {
				sb.WriteRune(rune(tile))
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (game Game) beam(i, j int, dir Direction) {
	if i < 0 || i >= game.h() || j < 0 || j >= game.w() {
		return
	}
	if game.isEnergized[i][j][dir] {
		return
	}
	game.isEnergized[i][j][dir] = true
	switch game.field[i][j] {
	case Empty:
		switch dir {
		case Up:
			game.beam(i-1, j, Up)
		case Left:
			game.beam(i, j-1, Left)
		case Right:
			game.beam(i, j+1, Right)
		case Down:
			game.beam(i+1, j, Down)
		}
	case MirrirSlash:
		switch dir {
		case Up:
			game.beam(i, j+1, Right)
		case Left:
			game.beam(i+1, j, Down)
		case Right:
			game.beam(i-1, j, Up)
		case Down:
			game.beam(i, j-1, Left)
		}
	case MirrorBack:
		switch dir {
		case Up:
			game.beam(i, j-1, Left)
		case Left:
			game.beam(i-1, j, Up)
		case Right:
			game.beam(i+1, j, Down)
		case Down:
			game.beam(i, j+1, Right)
		}
	case SplitHor:
		switch dir {
		case Up, Down:
			game.beam(i, j-1, Left)
			game.beam(i, j+1, Right)
		case Left:
			game.beam(i, j-1, Left)
		case Right:
			game.beam(i, j+1, Right)
		}
	case SplitVer:
		switch dir {
		case Left, Right:
			game.beam(i-1, j, Up)
			game.beam(i+1, j, Down)
		case Up:
			game.beam(i-1, j, Up)
		case Down:
			game.beam(i+1, j, Down)
		}
	}
}

func (game Game) clear() {
	for _, row := range game.isEnergized {
		for j := range row {
			row[j] = [4]bool{}
		}
	}
}

func (game Game) countEnergized() (count int) {
	for _, row := range game.isEnergized {
		for _, energized := range row {
			for _, dirEnergized := range energized {
				if dirEnergized {
					count++
					break
				}
			}
		}
	}
	return
}

func (game Game) maxCountEnergized() (maxCount int) {
	checkForBeam := func(i, j int, dir Direction) {
		game.beam(i, j, dir)
		count := game.countEnergized()
		if count > maxCount {
			fmt.Print(game)
			fmt.Println("New record:", i, j, dir, count)
			maxCount = count
		}
		game.clear()
	}
	for i := 0; i < game.h(); i++ {
		checkForBeam(i, 0, Right)
		checkForBeam(i, game.w()-1, Left)
	}
	for j := 0; j < game.w(); j++ {
		checkForBeam(0, j, Down)
		checkForBeam(game.h()-1, j, Up)
	}
	return maxCount
}

func parseRow(line string) (row []Tile) {
	row = make([]Tile, len(line))
	for i, r := range line {
		row[i] = Tile(r)
	}
	return
}

func appendRows(game Game, row []Tile) Game {
	game.field = append(game.field, row)
	game.isEnergized = append(game.isEnergized, make([]Energized, len(row)))
	return game
}

func Run() int {
	game := utils.ProcessInput("day16.txt", Game{}, parseRow, appendRows)
	return game.maxCountEnergized()
}
