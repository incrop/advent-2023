#!/bin/bash

set -e

DAY="$1"

mkdir "${DAY}"

cat > "${DAY}/${DAY}.go" << EOF
package ${DAY}

import "advent/utils"


func parseLine(line string) string {
	return line
}

func aggregate(acc int, elem string) int {
	return acc + len(elem)
}

func Run() int {
	utils.ProcessInput("${DAY}_test.txt", 0, parseLine, aggregate)
	return 0
}
EOF

cat > "advent.go" << EOF
package main

import (
	"advent/${DAY}"
	"fmt"
)

func main() {
	fmt.Println(${DAY}.Run())
}
EOF

touch "input/${DAY}_test.txt" "input/${DAY}.txt"