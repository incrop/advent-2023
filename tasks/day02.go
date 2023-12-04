package tasks

import (
	"advent/utils"
	"strconv"
	"strings"
)

type d02_Cubes struct {
	red   int
	green int
	blue  int
}

type d02_Game struct {
	id       int
	cubesets []d02_Cubes
}

func d02_parseCubes(str string) d02_Cubes {
	cubes := d02_Cubes{}
	for _, cubeStr := range utils.Fields(str, ",") {
		numAndColor := strings.Fields(cubeStr)
		num, _ := strconv.Atoi(numAndColor[0])
		switch numAndColor[1] {
		case "red":
			cubes.red = num
		case "green":
			cubes.green = num
		case "blue":
			cubes.blue = num
		}
	}
	return cubes
}

func d02_parseGame(str string) d02_Game {
	idAndCubesets := utils.Fields(str, ":;")
	idStr, _ := strings.CutPrefix(idAndCubesets[0], "Game ")
	id, _ := strconv.Atoi(idStr)
	cubesets := make([]d02_Cubes, len(idAndCubesets)-1, len(idAndCubesets)-1)
	for i, cubesetsStr := range idAndCubesets[1:] {
		cubesets[i] = d02_parseCubes(cubesetsStr)
	}
	return d02_Game{
		id,
		cubesets,
	}
}

func d02_sumCorrectIds(acc int, game d02_Game) int {
	for _, cubeset := range game.cubesets {
		if cubeset.red > 12 || cubeset.green > 13 || cubeset.blue > 14 {
			return acc
		}
	}
	return acc + game.id
}

func d02_sumMinimalPowers(acc int, game d02_Game) int {
	minCubeset := d02_Cubes{}
	for _, cubeset := range game.cubesets {
		if minCubeset.red < cubeset.red {
			minCubeset.red = cubeset.red
		}
		if minCubeset.green < cubeset.green {
			minCubeset.green = cubeset.green
		}
		if minCubeset.blue < cubeset.blue {
			minCubeset.blue = cubeset.blue
		}
	}
	minimalPower := minCubeset.red * minCubeset.green * minCubeset.blue
	return acc + minimalPower
}

func Day02() int {
	return utils.ProcessInput("day02.txt", 0, d02_parseGame, d02_sumMinimalPowers)
}
