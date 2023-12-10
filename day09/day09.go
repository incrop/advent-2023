package day09

import (
	"advent/utils"
	"fmt"
	"strconv"
	"strings"
)

type Row []int

func (row Row) isAllZeros() bool {
	for _, value := range row {
		if value != 0 {
			return false
		}
	}
	return true
}

type History []Row

func (history History) extrapolateForward() {
	increment := 0
	for i := len(history) - 1; i >= 0; i-- {
		row := history[i]
		increment += row[len(row)-1]
		history[i] = append(row, increment)
	}
}

func (history History) extrapolateBack() {
	decrement := 0
	for i := len(history) - 1; i >= 0; i-- {
		row := append([]int{0}, history[i]...)
		decrement = row[1] - decrement
		row[0] = decrement
		history[i] = row
	}
}

func (history History) String() string {
	var sb strings.Builder
	for i, row := range history {
		for j := 0; j < i; j++ {
			sb.WriteString("  ")
		}
		fmt.Fprintf(&sb, "%v\n", row)
	}
	return sb.String()
}

func (history History) lastValue() int {
	values := history[0]
	return values[len(values)-1]
}

func parseLine(line string) (history History) {
	fields := strings.Fields(line)
	initial := make(Row, len(fields), 2*len(fields))
	for i, field := range fields {
		value, _ := strconv.Atoi(field)
		initial[i] = value
	}
	history = make(History, 1)
	history[0] = initial
	row := initial
	for !row.isAllZeros() {
		nextRow := make(Row, len(row)-1, cap(row)-1)
		for i := range nextRow {
			nextRow[i] = row[i+1] - row[i]
		}
		history = append(history, nextRow)
		row = nextRow
	}
	return
}

func sumExtrapolatedForwardValues(acc int, history History) int {
	fmt.Println(history)
	history.extrapolateForward()
	fmt.Println(history)
	return acc + history.lastValue()
}

func sumExtrapolatedBackValues(acc int, history History) int {
	history.extrapolateBack()
	return acc + history[0][0]
}

func Run() int {
	//return utils.ProcessInput("day09.txt", 0, parseLine, sumExtrapolatedForwardValues)
	return utils.ProcessInput("day09.txt", 0, parseLine, sumExtrapolatedBackValues)
}
