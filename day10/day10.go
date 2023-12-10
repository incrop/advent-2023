package day10

import (
	"advent/utils"
	"fmt"
	"strings"
)

type Direction byte

const (
	Up Direction = iota
	Left
	Right
	Down
)

func (direction Direction) opposite() Direction {
	switch direction {
	case Up:
		return Down
	case Left:
		return Right
	case Right:
		return Left
	case Down:
		return Up
	default:
		return direction
	}
}

type Pipe [2]Direction

var StartingPosition = Pipe{Up, Up}
var Ground = Pipe{Down, Down}

func parsePipe(r rune) Pipe {
	switch r {
	case '|':
		return Pipe{Up, Down}
	case '-':
		return Pipe{Left, Right}
	case 'L':
		return Pipe{Up, Right}
	case 'J':
		return Pipe{Up, Left}
	case '7':
		return Pipe{Left, Down}
	case 'F':
		return Pipe{Right, Down}
	case '.':
		return Ground
	case 'S':
		return StartingPosition
	default:
		panic("Unknown pipe")
	}
}

func (pipe Pipe) isConnected(direction Direction) bool {
	if pipe[0] == pipe[1] {
		return false
	}
	return pipe[0] == direction || pipe[1] == direction
}

type Position struct {
	i, j int
}

type Color byte

const (
	NoColor Color = iota
	LoopColor
	InsideColor
)

type Tile struct {
	pipe  Pipe
	color Color
}

type Labyrinth struct {
	tiles    [][]Tile
	position Position
}

func (labyrinth Labyrinth) String() string {
	var sb strings.Builder
	for _, row := range labyrinth.tiles {
		for _, tile := range row {
			switch tile.color {
			case LoopColor:
				sb.WriteString("\033[32m")
			case InsideColor:
				sb.WriteString("\033[31m")
			}
			switch tile.pipe {
			case Pipe{Up, Down}:
				sb.WriteRune('║')
			case Pipe{Left, Right}:
				sb.WriteRune('═')
			case Pipe{Up, Right}:
				sb.WriteRune('╚')
			case Pipe{Up, Left}:
				sb.WriteRune('╝')
			case Pipe{Left, Down}:
				sb.WriteRune('╗')
			case Pipe{Right, Down}:
				sb.WriteRune('╔')
			case Ground:
				sb.WriteRune('.')
			default:
				sb.WriteRune('?')
			}
			if tile.color > 0 {
				sb.WriteString("\033[0m")
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (labyrinth Labyrinth) LoopArea() string {
	var sb strings.Builder
	for _, row := range labyrinth.tiles {
		isInside := false
		insideDirection := Right
		for _, tile := range row {
			switch tile.color {
			case NoColor:
				sb.WriteRune(' ')
			case LoopColor:
				sb.WriteString("\033[32m")
				switch tile.pipe {
				case Pipe{Up, Down}:
					if isInside {
						sb.WriteRune('▌')
					} else {
						sb.WriteRune('▐')
					}
					isInside = !isInside
				case Pipe{Left, Right}:
					if insideDirection == Up {
						sb.WriteRune('▀')
					} else {
						sb.WriteRune('▄')
					}
				case Pipe{Up, Right}:
					if isInside {
						sb.WriteRune('▙')
						insideDirection = Down
					} else {
						sb.WriteRune('▝')
						insideDirection = Up
					}
				case Pipe{Up, Left}:
					if insideDirection == Up {
						sb.WriteRune('▘')
						isInside = false
					} else {
						sb.WriteRune('▟')
						isInside = true
					}
				case Pipe{Left, Down}:
					if insideDirection == Up {
						sb.WriteRune('▜')
						isInside = true
					} else {
						sb.WriteRune('▖')
						isInside = false
					}
				case Pipe{Right, Down}:
					if isInside {
						sb.WriteRune('▛')
						insideDirection = Up
					} else {
						sb.WriteRune('▗')
						insideDirection = Down
					}
				default:
					sb.WriteRune('?')
				}
				sb.WriteString("\033[0m")
			case InsideColor:
				sb.WriteString("\033[31m█\033[0m")
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (labyrinth Labyrinth) colorMainLoop() (length int) {
	i := labyrinth.position.i
	j := labyrinth.position.j
	direction := labyrinth.tiles[i][j].pipe[0].opposite()
	for {
		labyrinth.tiles[i][j].color = LoopColor
		pipe := labyrinth.tiles[i][j].pipe
		if pipe[0].opposite() != direction {
			direction = pipe[0]
		} else {
			direction = pipe[1]
		}
		switch direction {
		case Up:
			i--
		case Left:
			j--
		case Right:
			j++
		case Down:
			i++
		}
		length++
		if i == labyrinth.position.i && j == labyrinth.position.j {
			return
		}
	}
}

func (labyrinth Labyrinth) colorInsideTiles() (num int) {
	for _, row := range labyrinth.tiles[1 : len(labyrinth.tiles)-1] {
		isInside := false
		lastUpOrDown := Right
		for j, tile := range row {
			if tile.color == LoopColor {
				switch tile.pipe {
				case Pipe{Up, Down}:
					isInside = !isInside
				case Pipe{Up, Right}:
					lastUpOrDown = Up
				case Pipe{Right, Down}:
					lastUpOrDown = Down
				case Pipe{Left, Right}:
				case Pipe{Up, Left}:
					if lastUpOrDown == Down {
						isInside = !isInside
					}
				case Pipe{Left, Down}:
					if lastUpOrDown == Up {
						isInside = !isInside
					}
				default:
					panic("Unexpected pipe")
				}
			} else if isInside {
				row[j].color = InsideColor
				num++
			}
		}
	}
	return
}

func (labyrinth Labyrinth) inferStartPositionPipe() {
	i := labyrinth.position.i
	j := labyrinth.position.j
	tiles := labyrinth.tiles
	idx := 0
	pipe := Pipe{}
	if i > 0 && tiles[i-1][j].pipe.isConnected(Down) {
		pipe[idx] = Up
		idx++
	}
	if j > 0 && tiles[i][j-1].pipe.isConnected(Right) {
		pipe[idx] = Left
		idx++
	}
	if j < len(tiles[i])-1 && tiles[i][j+1].pipe.isConnected(Left) {
		pipe[idx] = Right
		idx++
	}
	if idx < 2 {
		pipe[idx] = Down
	}
	tiles[i][j].pipe = pipe
}

func parseLine(line string) (row []Tile) {
	row = make([]Tile, len(line))
	for i, r := range line {
		row[i].pipe = parsePipe(r)
	}
	return
}

func aggregate(labyrinth Labyrinth, row []Tile) Labyrinth {
	i := len(labyrinth.tiles)
	labyrinth.tiles = append(labyrinth.tiles, row)
	for j, tile := range row {
		if tile.pipe == StartingPosition {
			labyrinth.position = Position{i, j}
			labyrinth.inferStartPositionPipe()
			break
		}
	}
	return labyrinth
}

func Run() int {
	labyrinth := utils.ProcessInput("day10.txt", Labyrinth{}, parseLine, aggregate)
	defer func() {
		fmt.Println(labyrinth)
		fmt.Println(labyrinth.LoopArea())
	}()
	labyrinth.colorMainLoop()
	return labyrinth.colorInsideTiles()
}
