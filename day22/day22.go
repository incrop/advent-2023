package day22

import (
	"advent/utils"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Coords struct {
	x, y, z int
}

type Brick [2]Coords

func (brick Brick) placeOnTop(z int) Brick {
	brick[1].z -= brick[0].z - z - 1
	brick[0].z = z + 1
	return brick
}

func (brick1 Brick) overlapsXY(brick2 Brick) bool {
	noOverlap := func(a1, a2, b1, b2 int) bool {
		if a2 < a1 {
			a1, a2 = a2, a1
		}
		if b2 < b1 {
			b1, b2 = b2, b1
		}
		return a2 < b1 || b2 < a1
	}
	if noOverlap(brick1[0].x, brick1[1].x, brick2[0].x, brick2[1].x) {
		return false
	}
	if noOverlap(brick1[0].y, brick1[1].y, brick2[0].y, brick2[1].y) {
		return false
	}
	return true
}

type PlacedBrick struct {
	id          int
	coords      Brick
	supportedBy []int
}

type System []PlacedBrick

func (system System) String() string {
	var sb strings.Builder
	for _, brick := range system {
		fmt.Fprintln(&sb, brick)
	}
	return sb.String()
}

func (system System) countOptionalBricks() int {
	mandatoryBricks := make([]bool, len(system))
	for _, brick := range system {
		if len(brick.supportedBy) == 1 {
			mandatoryBricks[brick.supportedBy[0]] = true
		}
	}
	return len(system) - len(mandatoryBricks)
}

func (system System) countFallen(id int) int {
	fallen := map[int]bool{id: true}
	for id := id + 1; id < len(system); id++ {
		brick := system[id]
		if brick.supportedBy == nil {
			continue
		}
		hasSupport := false
		for _, supportId := range brick.supportedBy {
			if !fallen[supportId] {
				hasSupport = true
				break
			}
		}
		if !hasSupport {
			fallen[id] = true
		}
	}
	return len(fallen) - 1
}

func (system System) sumOfFallingBricks() (sum int) {
	for i := range system {
		sum += system.countFallen(i)
	}
	return
}

func placeBricks(bricks []Brick) (system System) {
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i][0].z < bricks[j][0].z
	})
	for i, brick := range bricks {
		maxZ := 0
		var supportIds []int
		for j, placedBrick := range system {
			if !brick.overlapsXY(placedBrick.coords) {
				continue
			}
			placedBrickMaxZ := placedBrick.coords[1].z
			if placedBrickMaxZ > maxZ {
				maxZ = placedBrickMaxZ
				supportIds = []int{j}
			} else if placedBrickMaxZ == maxZ {
				supportIds = append(supportIds, j)
			}
		}
		system = append(system, PlacedBrick{
			id:          i,
			coords:      brick.placeOnTop(maxZ),
			supportedBy: supportIds,
		})
	}
	return
}

func parseBrick(line string) (brick Brick) {
	for i, coordStr := range utils.Fields(line, "~") {
		coord := &brick[i]
		for j, numStr := range utils.Fields(coordStr, ",") {
			num, _ := strconv.Atoi(numStr)
			switch j {
			case 0:
				coord.x = num
			case 1:
				coord.y = num
			case 2:
				coord.z = num
			}
		}
	}
	if brick[1].z < brick[0].z {
		brick[0], brick[1] = brick[1], brick[0]
	}
	return
}

func appendBrick(bricks []Brick, brick Brick) []Brick {
	return append(bricks, brick)
}

func Run() int {
	bricks := utils.ProcessInput("day22.txt", nil, parseBrick, appendBrick)
	system := placeBricks(bricks)
	return system.sumOfFallingBricks()
}
