package day14_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/tastapod/advent2022/day14"
	. "github.com/tastapod/advent2022/pair"
	. "github.com/tastapod/advent2022/segment"
	"strings"
	"testing"
)

func TestCreatesPairs(t *testing.T) {
	path := Path{{10, 20}, {10, 30}, {20, 30}, {40, 30}, {40, 80}}

	expected := []Pair[Point]{
		{Point{X: 10, Y: 20}, Point{X: 10, Y: 30}},
		{Point{X: 10, Y: 30}, Point{X: 20, Y: 30}},
		{Point{X: 20, Y: 30}, Point{X: 40, Y: 30}},
		{Point{X: 40, Y: 30}, Point{X: 40, Y: 80}},
	}

	pairs := ZipWithNext(path)
	assert.Equal(t, expected, pairs)
}

func TestExpandsPath(t *testing.T) {
	assert := assert.New(t)

	om := ObstacleMap{}
	path := Path{{498, 4}, {498, 6}, {496, 6}}

	err := om.PlotObstaclePath(path, ROCK)
	if err != nil {
		t.Error(err)
	}

	for _, point := range []Point{{498, 4}, {498, 5}, {498, 6}, {497, 6}, {496, 6}} {
		assert.Contains(om, point)
	}

	for _, point := range []Point{{497, 4}, {500, 5}, {498, 10}} {
		assert.NotContains(om, point)
	}
}

func TestParsesPath(t *testing.T) {
	input := "503,4 -> 502,4 -> 502,9 -> 494,9"
	expected := Path{{503, 4}, {502, 4}, {502, 9}, {494, 9}}

	points := ParsePath(input)
	assert.Equal(t, expected, points)
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
	assert.Equal(t, expected, paths)
}

func sampleCave() (Cave, error) {
	return NewCaveFromStrings(sampleInput)
}

func TestDropsSand(t *testing.T) {
	assert := assert.New(t)
	cave, _ := sampleCave()

	// drop first grain
	landedAt, landed := cave.DropSandFrom(StartPoint)
	assert.Equal(true, landed)
	assert.Equal(Point{X: 500, Y: 8}, landedAt)

	// then the next lot
	for i := 2; i <= 23; i++ {
		cave.DropSandFrom(StartPoint)
	}

	// this one should settle
	_, landed = cave.DropSandFrom(StartPoint)
	assert.Equal(true, landed)

	// this one should fall through
	_, landed = cave.DropSandFrom(StartPoint)
	assert.Equal(false, landed)
}

func TestCountsSand(t *testing.T) {
	cave, _ := sampleCave()
	numGrains := cave.FillWithSandFrom(StartPoint)
	assert.Equal(t, 24, numGrains)
}

func TestCalculatesBaseline(t *testing.T) {
	assert := assert.New(t)

	cave, _ := sampleCave()
	start, end := cave.BaselineAround(StartPoint)
	assert.Equal(Point{X: 500 - 12, Y: 11}, start)
	assert.Equal(Point{X: 500 + 12, Y: 11}, end)
}

func TestFillsWithBaseline(t *testing.T) {
	cave, _ := NewCaveFromStrings(sampleInput)
	err := cave.AddBaselineAround(StartPoint)
	if err != nil {
		t.Error(err)
	}
	numGrains := cave.FillWithSandFrom(StartPoint)
	assert.Equal(t, 93, numGrains)
}
