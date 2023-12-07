#!/bin/bash

set -e

DAY="$1"

mkdir "${DAY}"

cat > "${DAY}/${DAY}.go" << EOF
package ${DAY}

import "advent/utils"

func Run() int {
	utils.ProcessInput("${DAY}_test.txt", 0, utils.Identity, utils.JustPrint)
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