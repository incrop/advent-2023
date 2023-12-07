package day03

import (
	"advent/utils"
)

const (
	Empty  = -1
	Gear   = -2
	Symbol = -2
)

type Scheme struct {
	data    [][]int
	numbers []int
}

func (scheme Scheme) partNumbers() []int {
	isPartNum := make([]bool, len(scheme.numbers))
	for i, row := range scheme.data {
		for j, val := range row {
			if val != Gear && val != Symbol {
				continue
			}
			for m := max(i-1, 0); m <= min(i+1, len(scheme.data)-1); m++ {
				for n := max(j-1, 0); n <= min(j+1, len(row)-1); n++ {
					numIdx := scheme.data[m][n]
					if numIdx >= 0 {
						isPartNum[numIdx] = true
					}
				}
			}
		}
	}
	partNums := make([]int, 0, len(scheme.numbers))
	for i, num := range scheme.numbers {
		if isPartNum[i] {
			partNums = append(partNums, num)
		}
	}
	return partNums
}

func (scheme Scheme) gearRatios() []int {
	gearRatios := make([]int, 0)
	for i, row := range scheme.data {
		for j, val := range row {
			if val != Gear {
				continue
			}
			partNumbers := map[int]int{}
			for m := max(i-1, 0); m <= min(i+1, len(scheme.data)-1); m++ {
				for n := max(j-1, 0); n <= min(j+1, len(row)-1); n++ {
					numIdx := scheme.data[m][n]
					if numIdx >= 0 {
						partNumbers[numIdx] = scheme.numbers[numIdx]
					}
				}
			}
			if len(partNumbers) != 2 {
				continue
			}
			ratio := 1
			for _, num := range partNumbers {
				ratio *= num
			}
			gearRatios = append(gearRatios, ratio)
		}
	}
	return gearRatios
}

func collectScheme(scheme Scheme, row string) Scheme {
	dataRow := make([]int, len(row))
	numIdx := -1
	for i, c := range row {
		if c >= '0' && c <= '9' {
			if numIdx == -1 {
				scheme.numbers = append(scheme.numbers, 0)
				numIdx = len(scheme.numbers) - 1
			}
			scheme.numbers[numIdx] = scheme.numbers[numIdx]*10 + int(c-'0')
			dataRow[i] = len(scheme.numbers) - 1
		} else {
			numIdx = -1
			switch c {
			case '.':
				dataRow[i] = Empty
			case '*':
				dataRow[i] = Gear
			default:
				dataRow[i] = Symbol
			}
		}
	}
	scheme.data = append(scheme.data, dataRow)
	return scheme
}

func Run() int {
	scheme := utils.ProcessInput("day03.txt", Scheme{}, utils.Identity, collectScheme)
	sum := 0
	for _, num := range scheme.gearRatios() {
		sum += num
	}
	return sum
}
