package day14

import (
	"advent/utils"
	"math/big"
	"strings"
)

type Tile rune

const (
	Round Tile = 'O'
	Cube  Tile = '#'
	Empty Tile = '.'
)

type Platform [][]Tile

func (from Platform) clone() (to Platform) {
	to = make(Platform, len(from))
	for i, row := range from {
		to[i] = make([]Tile, len(row))
		copy(to[i], row)
	}
	return
}

func (platform Platform) tiltNorth() {
	h, w := len(platform), len(platform[0])
	for i1 := 1; i1 < h; i1++ {
		for j := 0; j < w; j++ {
			if platform[i1][j] == Round && platform[i1-1][j] == Empty {
				i2 := i1 - 1
				for i2 > 0 && platform[i2-1][j] == Empty {
					i2--
				}
				platform[i1][j], platform[i2][j] = platform[i2][j], platform[i1][j]
			}
		}
	}
}

func (platform Platform) tiltWest() {
	h, w := len(platform), len(platform[0])
	for j1 := 1; j1 < w; j1++ {
		for i := 0; i < h; i++ {
			if platform[i][j1] == Round && platform[i][j1-1] == Empty {
				j2 := j1 - 1
				for j2 > 0 && platform[i][j2-1] == Empty {
					j2--
				}
				platform[i][j1], platform[i][j2] = platform[i][j2], platform[i][j1]
			}
		}
	}
}

func (platform Platform) tiltSouth() {
	h, w := len(platform), len(platform[0])
	for i1 := h - 2; i1 >= 0; i1-- {
		for j := 0; j < w; j++ {
			if platform[i1][j] == Round && platform[i1+1][j] == Empty {
				i2 := i1 + 1
				for i2 < h-1 && platform[i2+1][j] == Empty {
					i2++
				}
				platform[i1][j], platform[i2][j] = platform[i2][j], platform[i1][j]
			}
		}
	}
}

func (platform Platform) tiltEast() {
	h, w := len(platform), len(platform[0])
	for j1 := w - 2; j1 >= 0; j1-- {
		for i := 0; i < h; i++ {
			if platform[i][j1] == Round && platform[i][j1+1] == Empty {
				j2 := j1 + 1
				for j2 < w-1 && platform[i][j2+1] == Empty {
					j2++
				}
				platform[i][j1], platform[i][j2] = platform[i][j2], platform[i][j1]
			}
		}
	}
}

func (platform Platform) fingerprint() string {
	var bits big.Int
	i := 0
	for _, row := range platform {
		for _, tile := range row {
			if tile == Empty {
				i++
			} else if tile == Round {
				bits.SetBit(&bits, i, 1)
				i++
			}
		}
	}
	return bits.String()
}

func (platform Platform) tiltCounterclockwise(times int) Platform {
	platform = platform.clone()
	fingerprintToIndex := make(map[string]int, 0)
	fingerprintToIndex[platform.fingerprint()] = 0
	states := make([]Platform, 0)
	states = append(states, platform.clone())
	for i := 1; i <= times; i++ {
		platform.tiltNorth()
		platform.tiltWest()
		platform.tiltSouth()
		platform.tiltEast()
		fingerprint := platform.fingerprint()
		if index, found := fingerprintToIndex[fingerprint]; found {
			loopLen := i - index
			iterationsLeft := times - i
			modulo := iterationsLeft % loopLen
			return states[index+modulo]
		}
		fingerprintToIndex[fingerprint] = len(states)
		states = append(states, platform.clone())
	}
	return platform
}

func (platform Platform) totalLoad() (load int) {
	for i, factor := len(platform)-1, 1; i >= 0; i, factor = i-1, factor+1 {
		for _, tile := range platform[i] {
			if tile == Round {
				load += factor
			}
		}
	}
	return
}

func (platform Platform) String() string {
	var sb strings.Builder
	for _, row := range platform {
		for _, land := range row {
			sb.WriteRune(rune(land))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func parseRow(line string) []Tile {
	row := make([]Tile, len(line))
	for i, r := range line {
		row[i] = Tile(r)
	}
	return row
}

func appendRow(platform Platform, row []Tile) Platform {
	return append(platform, row)
}

func Run() int {
	platform := utils.ProcessInput("day14.txt", make(Platform, 0), parseRow, appendRow)
	platform = platform.tiltCounterclockwise(1000000000)
	return platform.totalLoad()
}
