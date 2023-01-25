package day14_test

import (
	"github.com/tastapod/advent2022/check"
	. "github.com/tastapod/advent2022/day14"
	"strings"
	"testing"
)

func TestCreatesPairs(t *testing.T) {
	var expected = []Segment{
		{Point{X: 10, Y: 20}, Point{X: 25, Y: 35}},
		{Point{X: 25, Y: 35}, Point{X: 40, Y: 50}},
		{Point{X: 40, Y: 50}, Point{X: 55, Y: 65}},
		{Point{X: 55, Y: 65}, Point{X: 70, Y: 80}},
	}

	pairs := ZipWithNext([]Point{{10, 20}, {25, 35}, {40, 50}, {55, 65}, {70, 80}})

	check.Equal(t, expected, pairs)
}

func TestExpandsPath(t *testing.T) {
	cave := NewCave()
	cave.ExpandPath([]Point{{498, 4}, {498, 6}, {496, 6}})

	for _, point := range []Point{{498, 4}, {498, 5}, {498, 6}, {497, 6}, {496, 6}} {
		if !cave.ObstacleAt(point) {
			t.Error(point, "expected but not found")
		}
	}

	for _, point := range []Point{{497, 4}, {500, 5}, {498, 10}} {
		if cave.ObstacleAt(point) {
			t.Error(point, "found but not expected")
		}
	}
}

func TestParsesPath(t *testing.T) {
	input := "503,4 -> 502,4 -> 502,9 -> 494,9"
	expected := Path{{503, 4}, {502, 4}, {502, 9}, {494, 9}}

	points := ParsePath(input)
	check.Equal(t, expected, points)
}

var sampleInput = strings.Split(strings.TrimSpace(`
498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`), "\n")

func TestParsesMultiplePaths(t *testing.T) {

	expected := []Path{
		{{498, 4}, {498, 6}, {496, 6}},
		{{503, 4}, {502, 4}, {502, 9}, {494, 9}},
	}

	paths := ParsePaths(sampleInput)
	check.Equal(t, expected, paths)
}

func sampleCave() Cave {
	return NewCaveFromStrings(sampleInput)
}

func TestDropsSand(t *testing.T) {
	cave := sampleCave()

	// drop first grain
	landedAt, landed := cave.DropSandFrom(StartPoint)
	check.Equal(t, true, landed)
	check.Equal(t, Point{X: 500, Y: 8}, landedAt)

	// then the next lot
	for i := 2; i <= 23; i++ {
		cave.DropSandFrom(StartPoint)
	}

	// this one should settle
	_, landed = cave.DropSandFrom(StartPoint)
	check.Equal(t, true, landed)

	// this one should fall through
	_, landed = cave.DropSandFrom(StartPoint)
	check.Equal(t, false, landed)
}

func TestCountsSand(t *testing.T) {
	cave := sampleCave()
	numGrains := cave.FillWithSandFrom(StartPoint)
	check.Equal(t, 24, numGrains)
}

func TestCalculatesBaseline(t *testing.T) {
	cave := sampleCave()
	start, end := cave.BaselineAround(StartPoint)
	check.Equal(t, Point{X: 500 - 12, Y: 11}, start)
	check.Equal(t, Point{X: 500 + 12, Y: 11}, end)
}

func TestFillsWithBaseline(t *testing.T) {
	cave := NewCaveFromStrings(sampleInput)
	cave.AddBaselineAround(StartPoint)
	numGrains := cave.FillWithSandFrom(StartPoint)
	check.Equal(t, 93, numGrains)
}
