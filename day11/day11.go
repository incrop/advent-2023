package day11

import (
	"advent/utils"
	"sort"
)

type Galaxy struct {
	x, y int
}

func (galaxy Galaxy) distanceTo(other Galaxy) int {
	return utils.Abs(galaxy.x-other.x) + utils.Abs(galaxy.y-other.y)
}

type Galaxies []Galaxy

func (galaxies Galaxies) expand(factor int) {
	if galaxies == nil {
		return
	}
	expandDimension := func(get func(Galaxy) int, set func(*Galaxy, int)) {
		sort.Slice(galaxies, func(i, j int) bool {
			return get(galaxies[i]) < get(galaxies[j])
		})
		totalDelta := 0
		lastVal := get(galaxies[0])
		for i := 1; i < len(galaxies); i++ {
			val := get(galaxies[i])
			if delta := val - lastVal; delta > 1 {
				totalDelta += (delta - 1) * (factor - 1)
			}
			lastVal = val
			set(&galaxies[i], val+totalDelta)
		}
	}
	expandDimension(
		func(galaxy Galaxy) int {
			return galaxy.x
		},
		func(galaxy *Galaxy, x int) {
			galaxy.x = x
		},
	)
	expandDimension(
		func(galaxy Galaxy) int {
			return galaxy.y
		},
		func(galaxy *Galaxy, y int) {
			galaxy.y = y
		},
	)
}

func (galaxies Galaxies) distancePairwiseSum() (sum int) {
	if len(galaxies) <= 1 {
		return
	}
	for len(galaxies) >= 2 {
		galaxy := galaxies[0]
		galaxies = galaxies[1:]
		for _, other := range galaxies {
			distance := galaxy.distanceTo(other)
			sum += distance
		}
	}
	return
}

func parseGalaxiesRow(line string, y int) (xs Galaxies) {
	xs = make(Galaxies, 0)
	for x, r := range line {
		if r == '#' {
			xs = append(xs, Galaxy{x: x, y: y})
		}
	}
	return
}

func appendAll(galaxies Galaxies, galaxiesRow Galaxies) Galaxies {
	return append(galaxies, galaxiesRow...)
}

func Run() int {
	galaxies := utils.ProcessInputWithLineNumbers("day11.txt", nil, parseGalaxiesRow, appendAll)
	galaxies.expand(1000000)
	return galaxies.distancePairwiseSum()
}
