package day14

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Segment struct {
	start, end Point
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
	return NewCaveFromPaths(parsePaths(paths))
}

func NewCaveFromPaths(paths []Path) Cave {
	result := NewCave()
	for _, path := range paths {
		result.expandPath(path)
	}
	return result
}

func NewCave() Cave {
	return Cave{make(ObstacleMap), 0}
}

// dropSandFrom drops a grain of sand from the start point until it settles.
// It returns the final resting place of the sand and true if the sand settles,
// or the start point and false if it falls through, past lowest point.
func (cave *Cave) dropSandFrom(start Point) (Point, bool) {
	sand := start
	if cave.isBlocked(sand) {
		return sand, false
	}

grains:
	for below := sand; sand.y <= cave.lowest; {
		for _, xBelow := range []int{sand.x, sand.x - 1, sand.x + 1} {
			below = Point{xBelow, sand.y + 1}
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

	for settled := true; settled; _, settled = cave.dropSandFrom(start) {
	}
	after := len(cave.obstacles)
	return after - before
}

func (cave *Cave) AddBaselineAround(centre Point) {
	left, right := cave.baselineAround(centre)
	cave.expandPath([]Point{left, right})
}

func (cave *Cave) expandPath(points Path) {
	segments := zipWithNext(points)

	for _, segment := range segments {
		cave.expandSegment(segment)
	}
}

func (cave *Cave) isBlocked(point Point) bool {
	_, exists := cave.obstacles[point]
	return exists
}

func (cave *Cave) expandSegment(segment Segment) {
	x1, x2 := inOrder(segment.start.x, segment.end.x)
	y1, y2 := inOrder(segment.start.y, segment.end.y)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			cave.obstacles[Point{x, y}] = ROCK
		}
	}
	if y2 > cave.lowest {
		cave.lowest = y2
	}
}

func (cave *Cave) baselineAround(start Point) (Point, Point) {
	depth := cave.lowest + 2
	return Point{start.x - depth - 1, depth}, Point{start.x + depth + 1, depth}
}

func zipWithNext(points Path) []Segment {
	ends := points[1:]
	result := make([]Segment, len(ends))

	for i, end := range ends {
		result[i] = Segment{points[i], end}
	}
	return result
}

func inOrder(i1 int, i2 int) (int, int) {
	if i1 <= i2 {
		return i1, i2
	} else {
		return i2, i1
	}
}

func parsePath(input string) Path {
	coords := strings.Split(input, " -> ")
	result := make([]Point, len(coords))

	for i, coord := range coords {
		x, y, _ := strings.Cut(coord, ",")
		result[i] = Point{toInt(x), toInt(y)}
	}
	return result
}

func parsePaths(input []string) []Path {
	result := make([]Path, len(input))

	for i, line := range input {
		result[i] = parsePath(line)
	}
	return result
}

func toInt(s string) int {
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}
	panic(fmt.Sprintf("'%s': cannot parse int", s))
}
