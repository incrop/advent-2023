package day13

import (
	"advent/utils"
	"fmt"
	"strings"
)

type Land rune

const (
	ASH  Land = '.'
	ROCK Land = '#'
)

type Pattern [][]Land

func (pattern Pattern) String() string {
	var sb strings.Builder
	for _, row := range pattern {
		for _, land := range row {
			sb.WriteRune(rune(land))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

const TARGET_SMUDGE_COUNT = 1

func (pattern Pattern) findVerticalReflection() int {
	h, w := len(pattern), len(pattern[0])
REFLECTION_INDEX:
	for index := 1; index < w; index++ {
		smudgeCount := 0
		for j1, j2 := index-1, index; j1 >= 0 && j2 < w; j1, j2 = j1-1, j2+1 {
			for i := 0; i < h; i++ {
				if pattern[i][j1] != pattern[i][j2] {
					if smudgeCount == TARGET_SMUDGE_COUNT {
						continue REFLECTION_INDEX
					} else {
						smudgeCount++
					}
				}
			}
		}
		if smudgeCount == TARGET_SMUDGE_COUNT {
			return index
		}
	}
	return 0
}

func (pattern Pattern) findHorizontalReflection() int {
	h, w := len(pattern), len(pattern[0])
REFLECTION_INDEX:
	for index := 1; index < h; index++ {
		smudgeCount := 0
		for i1, i2 := index-1, index; i1 >= 0 && i2 < h; i1, i2 = i1-1, i2+1 {
			for j := 0; j < w; j++ {
				if pattern[i1][j] != pattern[i2][j] {
					if smudgeCount == TARGET_SMUDGE_COUNT {
						continue REFLECTION_INDEX
					} else {
						smudgeCount++
					}
				}
			}
		}
		if smudgeCount == TARGET_SMUDGE_COUNT {
			return index
		}
	}
	return 0
}

func parseRow(line string) []Land {
	if line == "" {
		return nil
	}
	row := make([]Land, len(line))
	for i, r := range line {
		row[i] = Land(r)
	}
	return row
}

func appendRow(patterns []Pattern, row []Land) []Pattern {
	if row == nil {
		return append(patterns, nil)
	}
	i := len(patterns) - 1
	patterns[i] = append(patterns[i], row)
	return patterns
}

func Run() int {
	patterns := utils.ProcessInput(
		"day13.txt",
		make([]Pattern, 1),
		parseRow,
		appendRow,
	)
	total := 0
	for _, pattern := range patterns {
		summary := pattern.findVerticalReflection() + 100*pattern.findHorizontalReflection()
		if summary == 0 {
			fmt.Println(pattern)
		}
		total += summary
	}
	return total
}
