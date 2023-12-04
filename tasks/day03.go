package tasks

import (
	"advent/utils"
)

const (
	d03_Empty  = -1
	d03_Gear   = -2
	d03_Symbol = -2
)

type d03_Scheme struct {
	data    [][]int
	numbers []int
}

func (scheme d03_Scheme) partNumbers() []int {
	isPartNum := make([]bool, len(scheme.numbers))
	for i, row := range scheme.data {
		for j, val := range row {
			if val != d03_Gear && val != d03_Symbol {
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

func (scheme d03_Scheme) gearRatios() []int {
	gearRatios := make([]int, 0)
	for i, row := range scheme.data {
		for j, val := range row {
			if val != d03_Gear {
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

func d03_collectScheme(scheme d03_Scheme, row string) d03_Scheme {
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
				dataRow[i] = d03_Empty
			case '*':
				dataRow[i] = d03_Gear
			default:
				dataRow[i] = d03_Symbol
			}
		}
	}
	scheme.data = append(scheme.data, dataRow)
	return scheme
}

func Day03() int {
	scheme := utils.ProcessInput("day03.txt", d03_Scheme{}, utils.Identity, d03_collectScheme)
	sum := 0
	for _, num := range scheme.gearRatios() {
		sum += num
	}
	return sum
}
