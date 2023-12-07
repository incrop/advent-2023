package day02

import (
	"advent/utils"
	"strconv"
	"strings"
)

type Cubes struct {
	red   int
	green int
	blue  int
}

type Game struct {
	id       int
	cubesets []Cubes
}

func parseCubes(str string) Cubes {
	cubes := Cubes{}
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

func parseGame(str string) Game {
	idAndCubesets := utils.Fields(str, ":;")
	idStr, _ := strings.CutPrefix(idAndCubesets[0], "Game ")
	id, _ := strconv.Atoi(idStr)
	cubesets := make([]Cubes, len(idAndCubesets)-1, len(idAndCubesets)-1)
	for i, cubesetsStr := range idAndCubesets[1:] {
		cubesets[i] = parseCubes(cubesetsStr)
	}
	return Game{
		id,
		cubesets,
	}
}

func sumCorrectIds(acc int, game Game) int {
	for _, cubeset := range game.cubesets {
		if cubeset.red > 12 || cubeset.green > 13 || cubeset.blue > 14 {
			return acc
		}
	}
	return acc + game.id
}

func sumMinimalPowers(acc int, game Game) int {
	minCubeset := Cubes{}
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

func Run() int {
	return utils.ProcessInput("day02.txt", 0, parseGame, sumMinimalPowers)
}
