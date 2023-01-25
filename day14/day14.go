package day14

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

type Segment struct {
	Start, End Point
}

type Path []Point

type Obstacle rune

const (
	SAND Obstacle = 'o'
	ROCK Obstacle = '#'
)

type ObstacleMap map[Point]Obstacle

type Cave struct {
	obstacles ObstacleMap
	lowest    int
}

var StartPoint = Point{500, 0}

func NewCaveFromStrings(paths []string) Cave {
	return NewCaveFromPaths(ParsePaths(paths))
}

func NewCaveFromPaths(paths []Path) Cave {
	result := NewCave()
	for _, path := range paths {
		result.ExpandPath(path)
	}
	return result
}

func NewCave() Cave {
	return Cave{make(ObstacleMap), 0}
}

func (cave *Cave) ObstacleAt(point Point) bool {
	_, found := cave.obstacles[point]
	return found
}

// DropSandFrom drops a grain of sand from the start point until it settles.
// It returns the final resting place of the sand and true if the sand settles,
// or the start point and false if it falls through, past lowest point.
func (cave *Cave) DropSandFrom(start Point) (Point, bool) {
	sand := start
	if cave.isBlocked(sand) {
		return sand, false
	}

grains:
	for below := sand; sand.Y <= cave.lowest; {
		for _, xBelow := range []int{sand.X, sand.X - 1, sand.X + 1} {
			below = Point{xBelow, sand.Y + 1}
			if !cave.isBlocked(below) {
				sand = below
				continue grains
			}
		}
		// we have come to rest
		cave.obstacles[sand] = SAND
		return sand, true
	}
	// we fell through
	return start, false
}

func (cave *Cave) FillWithSandFrom(start Point) int {
	before := len(cave.obstacles)

	for settled := true; settled; _, settled = cave.DropSandFrom(start) {
	}
	after := len(cave.obstacles)
	return after - before
}

func (cave *Cave) AddBaselineAround(centre Point) {
	left, right := cave.BaselineAround(centre)
	cave.ExpandPath([]Point{left, right})
}

func (cave *Cave) ExpandPath(points Path) {
	segments := ZipWithNext(points)

	for _, segment := range segments {
		cave.expandSegment(segment)
	}
}

func (cave *Cave) isBlocked(point Point) bool {
	_, exists := cave.obstacles[point]
	return exists
}

func (cave *Cave) expandSegment(segment Segment) {
	x1, x2 := inOrder(segment.Start.X, segment.End.X)
	y1, y2 := inOrder(segment.Start.Y, segment.End.Y)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			cave.obstacles[Point{x, y}] = ROCK
		}
	}
	if y2 > cave.lowest {
		cave.lowest = y2
	}
}

func inOrder(i1 int, i2 int) (int, int) {
	if i1 <= i2 {
		return i1, i2
	} else {
		return i2, i1
	}
}

func (cave *Cave) BaselineAround(start Point) (Point, Point) {
	depth := cave.lowest + 2
	return Point{start.X - depth - 1, depth}, Point{start.X + depth + 1, depth}
}

func ZipWithNext(points Path) []Segment {
	ends := points[1:]
	result := make([]Segment, len(ends))

	for i, end := range ends {
		result[i] = Segment{points[i], end}
	}
	return result
}

func ParsePath(input string) Path {
	coords := strings.Split(input, " -> ")
	result := make([]Point, len(coords))

	for i, coord := range coords {
		x, y, _ := strings.Cut(coord, ",")
		result[i] = Point{toInt(x), toInt(y)}
	}
	return result
}

func ParsePaths(input []string) []Path {
	result := make([]Path, len(input))

	for i, line := range input {
		result[i] = ParsePath(line)
	}
	return result
}

func toInt(s string) int {
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}
	panic(fmt.Sprintf("'%s': cannot parse int", s))
}
