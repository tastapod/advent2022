package main

import (
	"fmt"
	"github.com/tastapod/advent2022/day14"
	"github.com/tastapod/advent2022/input"
	"strings"
)

func main() {
	paths := strings.Split(input.ForDay(14), "\n")
	cave := day14.NewCaveFromStrings(paths)
	part1 := cave.FillWithSand()

	fmt.Printf("Day 14 part 1: %d", part1)
}
