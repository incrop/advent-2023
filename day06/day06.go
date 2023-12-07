package day06

import (
	"advent/utils"
	"strings"
)

type Race struct {
	time     int
	distance int
}

func (race Race) marginOfError() (wins int) {
	for acceleration := 0; acceleration <= race.time; acceleration++ {
		distance := acceleration * (race.time - acceleration)
		if distance > race.distance {
			wins++
		}
	}
	return
}

type Races []Race

func (races Races) marginOfError() (product int) {
	product = 1
	for _, race := range races {
		product *= race.marginOfError()
	}
	return
}

const IGNORE_SPACES = true

func parseRaceLine(races Races, line string) Races {
	timesStr, found := strings.CutPrefix(line, "Time:")
	if IGNORE_SPACES {
		timesStr = strings.ReplaceAll(timesStr, " ", "")
	}
	if found {
		times := utils.ParseNumbers(timesStr)
		races = make([]Race, len(times))
		for i, time := range times {
			races[i].time = time
		}
		return races
	}
	distanceStr, found := strings.CutPrefix(line, "Distance:")
	if IGNORE_SPACES {
		distanceStr = strings.ReplaceAll(distanceStr, " ", "")
	}
	if found {
		for i, distance := range utils.ParseNumbers(distanceStr) {
			races[i].distance = distance
		}
		return races
	}
	panic("what is it?")
}

func Run() int {
	races := utils.ProcessInput("day06.txt", nil, utils.Identity, parseRaceLine)
	return races.marginOfError()
}
