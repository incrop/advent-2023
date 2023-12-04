package utils

import (
	"bufio"
	"log"
	"os"
)

func ProcessInput[T any, R any](name string, seed R, parseLine func(string) T, join func(R, T) R) R {
	file, err := os.Open("input/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	result := seed
	for scanner.Scan() {
		parsed := parseLine(scanner.Text())
		result = join(result, parsed)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func Sum(a int, b int) int {
	return a + b
}

func Identity[T any](a T) T {
	return a
}
