package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ProcessInput[T any, R any](name string, seed R, parseLine func(string) T, join func(R, T) R) R {
	parseLineIgnoringNumbers := func(line string, _ int) T {
		return parseLine(line)
	}
	return ProcessInputWithLineNumbers(name, seed, parseLineIgnoringNumbers, join)
}

func ProcessInputWithLineNumbers[T any, R any](name string, seed R, parseLine func(string, int) T, join func(R, T) R) R {
	file, err := os.Open("input/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	result := seed
	lineNumber := 0
	for scanner.Scan() {
		parsed := parseLine(scanner.Text(), lineNumber)
		result = join(result, parsed)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func Sum(a int, b int) int {
	return a + b
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Identity[T any](a T) T {
	return a
}

func JustPrint[R any, T any](r R, t T) R {
	fmt.Println(t)
	return r
}

func Fields(str string, spaces string) []string {
	f := func(c rune) bool {
		return strings.ContainsRune(spaces, c)
	}
	return strings.FieldsFunc(str, f)
}

func ParseNumbers(str string) []int {
	fields := strings.Fields(str)
	numbers := make([]int, len(fields))
	for i, field := range fields {
		number, err := strconv.Atoi(field)
		if err != nil {
			panic(err)
		}
		numbers[i] = number
	}
	return numbers
}
