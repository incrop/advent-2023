package day09

import "advent/utils"


func parseLine(line string) string {
	return line
}

func aggregate(acc int, elem string) int {
	return acc + len(elem)
}

func Run() int {
	utils.ProcessInput("day09_test.txt", 0, parseLine, aggregate)
	return 0
}
