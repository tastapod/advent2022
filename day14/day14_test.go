package day14

import (
	"github.com/tastapod/advent2022/check"
	"strings"
	"testing"
)

func TestCreatesPairs(t *testing.T) {
	var expected = []Segment{
		{Point{10, 20}, Point{25, 35}},
		{Point{25, 35}, Point{40, 50}},
		{Point{40, 50}, Point{55, 65}},
		{Point{55, 65}, Point{70, 80}},
	}

	pairs := zipWithNext([]Point{{10, 20}, {25, 35}, {40, 50}, {55, 65}, {70, 80}})

	check.Equal(t, expected, pairs)
}

func TestExpandsPath(t *testing.T) {
	cave := NewCave()
	cave.expandPath([]Point{{498, 4}, {498, 6}, {496, 6}})

	for _, point := range []Point{{498, 4}, {498, 5}, {498, 6}, {497, 6}, {496, 6}} {
		_, found := cave.obstacles[point]
		if !found {
			t.Error(point, "expected but not found")
		}
	}

	for _, point := range []Point{{497, 4}, {500, 5}, {498, 10}} {
		_, found := cave.obstacles[point]
		if found {
			t.Error(point, "found but not expected")
		}
	}
}

func TestParsesPath(t *testing.T) {
	input := "503,4 -> 502,4 -> 502,9 -> 494,9"
	expected := Path{{503, 4}, {502, 4}, {502, 9}, {494, 9}}

	points := parsePath(input)
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

	paths := parsePaths(sampleInput)
	check.Equal(t, expected, paths)
}

func sampleCave() Cave {
	return NewCaveFromStrings(sampleInput)
}

func TestDropsSand(t *testing.T) {
	cave := sampleCave()

	// drop first grain
	landedAt, landed := cave.dropSandFrom(StartPoint)
	check.Equal(t, true, landed)
	check.Equal(t, Point{500, 8}, landedAt)

	// then the next lot
	for i := 2; i <= 23; i++ {
		cave.dropSandFrom(StartPoint)
	}

	// this one should settle
	_, landed = cave.dropSandFrom(StartPoint)
	check.Equal(t, true, landed)

	// this one should fall through
	_, landed = cave.dropSandFrom(StartPoint)
	check.Equal(t, false, landed)
}

func TestCountsSand(t *testing.T) {
	cave := sampleCave()
	numGrains := cave.FillWithSandFrom(StartPoint)
	check.Equal(t, 24, numGrains)
}

func TestCalculatesBaseline(t *testing.T) {
	cave := sampleCave()
	start, end := cave.baselineAround(StartPoint)
	check.Equal(t, Point{500 - 12, 11}, start)
	check.Equal(t, Point{500 + 12, 11}, end)
}

func TestFillsWithBaseline(t *testing.T) {
	cave := NewCaveFromStrings(sampleInput)
	cave.AddBaselineAround(StartPoint)
	numGrains := cave.FillWithSandFrom(StartPoint)
	check.Equal(t, 93, numGrains)
}
