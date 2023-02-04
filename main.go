package main

import (
	"fmt"
	"github.com/tastapod/advent2022/day14"
	"github.com/tastapod/advent2022/day15"
	"github.com/tastapod/advent2022/input"
	"strings"
	"time"
)

func main() {
	solveDay14()
	solveDay15()
}

func solveDay14() {
	paths := strings.Split(input.ForDay(14), "\n")
	cave, _ := day14.NewCaveFromStrings(paths)
	part1 := cave.FillWithSandFrom(day14.StartPoint)
	fmt.Printf("Day 14 part 1: %d\n", part1)

	cave, _ = day14.NewCaveFromStrings(paths)
	if err := cave.AddBaselineAround(day14.StartPoint); err != nil {
		fmt.Println("Day 14", err)
	}
	part2 := cave.FillWithSandFrom(day14.StartPoint)
	fmt.Printf("Day 14 part 2: %d\n", part2)
}

func solveDay15() {
	readings, err := day15.ParseReadings(strings.Split(input.ForDay(15), "\n"))
	if err != nil {
		return
	}
	part1 := day15.SumOverlappingSegmentLengths(2000000, readings)
	fmt.Printf("Day 15 part 1: %d\n", part1)

	start := time.Now()
	point, _ := day15.FindVacantPoint(4_000_000, readings)
	part2 := day15.TuningFrequency(point)
	elapsed := time.Since(start)
	fmt.Printf("Day 15 part 2: %v -> %d (took %v)\n", point, part2, elapsed)

}
