package day12

import (
	"advent/utils"
	"fmt"
	"strconv"
	"strings"
)

type Condition rune

const (
	Operational Condition = '.'
	Broken      Condition = '#'
	Unknown     Condition = '?'
)

type ConditionRecord struct {
	brokenSeries []int
	conditions   []Condition
}

func (record ConditionRecord) String() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for _, c := range record.conditions {
		sb.WriteRune(rune(c))
	}
	for _, s := range record.brokenSeries {
		sb.WriteRune(' ')
		fmt.Fprintf(&sb, "%d", s)
	}
	sb.WriteRune(']')
	return sb.String()
}

var cache map[string]int = make(map[string]int)

func (record ConditionRecord) numberOfArrangements() int {
	if cached, ok := cache[record.String()]; ok {
		return cached
	}
	number := calculateNumberOfArrangements(record.brokenSeries, record.conditions)
	cache[record.String()] = number
	return number
}

func calculateNumberOfArrangements(brokenSeries []int, conditions []Condition) int {
	if len(brokenSeries) == 0 {
		for _, c := range conditions {
			if c == Broken {
				return 0
			}
		}
		return 1
	}
	currentSeriesLen := brokenSeries[0]
	for {
		if len(conditions) < currentSeriesLen+1 {
			return 0
		}
		if conditions[0] != Operational {
			break
		}
		conditions = conditions[1:]
	}

	canBePlacedHere := true
	if conditions[currentSeriesLen] == Broken {
		canBePlacedHere = false
	} else {
		for _, c := range conditions[0:currentSeriesLen] {
			if c == Operational {
				canBePlacedHere = false
				break
			}
		}
	}
	canBePlacedElsewhere := conditions[0] != Broken

	number := 0
	if canBePlacedHere {
		number += ConditionRecord{
			brokenSeries: brokenSeries[1:],
			conditions:   conditions[currentSeriesLen+1:],
		}.numberOfArrangements()
	}
	if canBePlacedElsewhere {
		number += ConditionRecord{
			brokenSeries: brokenSeries,
			conditions:   conditions[1:],
		}.numberOfArrangements()
	}
	return number
}

func parseConditionRecord(line string) (record ConditionRecord) {
	split := strings.Fields(line)
	conditionsStr := split[0]
	record.conditions = make([]Condition, len(conditionsStr)+1)
	for i, r := range conditionsStr {
		record.conditions[i] = Condition(r)
	}
	record.conditions[len(conditionsStr)] = Operational
	brokenSeriesFields := utils.Fields(split[1], ",")
	record.brokenSeries = make([]int, len(brokenSeriesFields))
	for i, numStr := range brokenSeriesFields {
		record.brokenSeries[i], _ = strconv.Atoi(numStr)
	}
	return
}

func (record ConditionRecord) multiply(n int) (result ConditionRecord) {
	result.conditions = make([]Condition, 0, len(record.conditions)*n)
	result.brokenSeries = make([]int, 0, len(record.brokenSeries)*n)
	for i := 0; i < n; i++ {
		result.conditions = append(result.conditions, record.conditions...)
		if i < n-1 {
			result.conditions[len(result.conditions)-1] = Unknown
		}
		result.brokenSeries = append(result.brokenSeries, record.brokenSeries...)
	}
	return
}

func aggregate(acc int, record ConditionRecord) int {
	return acc + record.multiply(5).numberOfArrangements()
}

func Run() int {
	return utils.ProcessInput("day12.txt", 0, parseConditionRecord, aggregate)
}
