package main

import (
	"fmt"
	"github.com/tastapod/advent2022/day14"
	"github.com/tastapod/advent2022/input"
	"strings"
)

func main() {
	solveDay14()
}

func solveDay14() {
	paths := strings.Split(input.ForDay(14), "\n")
	cave := day14.NewCaveFromStrings(paths)
	part1 := cave.FillWithSandFrom(day14.StartPoint)
	fmt.Printf("Day 14 part 1: %d\n", part1)

	cave = day14.NewCaveFromStrings(paths)
	cave.AddBaselineAround(day14.StartPoint)
	part2 := cave.FillWithSandFrom(day14.StartPoint)
	fmt.Printf("Day 14 part 2: %d\n", part2)
}
