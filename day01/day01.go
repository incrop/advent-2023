package day01

import (
	"advent/utils"
	"strings"
)

var textToNumber = map[string]int{
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"0":     0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
	"zero":  0,
}

func extractNumber(str string) int {
	first_digit := 0
	last_digit := 0
	for pos := range str {
		substr := str[pos:]
		for text, number := range textToNumber {
			if strings.HasPrefix(substr, text) {
				last_digit = number
				if first_digit == 0 {
					first_digit = last_digit
				}
			}
		}
	}
	return first_digit*10 + last_digit
}

func Run() int {
	return utils.ProcessInput("day01.txt", 0, extractNumber, utils.Sum)
}
