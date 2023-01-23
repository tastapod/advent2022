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

type Cave struct {
	start     Point
	obstacles ObstacleMap
}

func NewCaveFromStrings(paths []string) Cave {
	return NewCaveFromPaths(parsePaths(paths))
}

func NewCaveFromPaths(paths []Path) Cave {
	cave := Cave{Point{500, 0}, NewObstacleMap()}
	for _, path := range paths {
		cave.obstacles.expandPath(path)
	}
	return cave
}

func (c *Cave) baseline() (Point, Point) {
	return c.obstacles.baselineAround(c.start)
}

func (c *Cave) dropSand() (Point, bool) {
	return c.obstacles.dropSandFrom(c.start)
}

func (c *Cave) FillWithSand() int {
	return c.obstacles.fillWithSandFrom(c.start)
}

func (c *Cave) AddBaseline() {
	c.obstacles.addBaselineAround(c.start)
}

type ObstacleMap struct {
	points map[Point]Obstacle
	lowest int
}

func NewObstacleMap() ObstacleMap {
	return ObstacleMap{make(map[Point]Obstacle), 0}
}

// dropSandFrom drops a grain of sand from the start point until it settles.
// It returns the final resting place of the sand and true if the sand settles,
// or the start point and false if it falls through, past lowest point.
func (om *ObstacleMap) dropSandFrom(start Point) (Point, bool) {
	sand := start
	if om.isBlocked(sand) {
		return sand, false
	}

grains:
	for below := sand; sand.y <= om.lowest; {
		for _, xBelow := range []int{sand.x, sand.x - 1, sand.x + 1} {
			below = Point{xBelow, sand.y + 1}
			if !om.isBlocked(below) {
				sand = below
				continue grains
			}
		}
		// we have come to rest
		om.points[sand] = SAND
		return sand, true
	}
	// we fell through
	return start, false
}

func (om *ObstacleMap) fillWithSandFrom(start Point) int {
	before := len(om.points)

	for settled := true; settled; _, settled = om.dropSandFrom(start) {
	}
	after := len(om.points)
	return after - before
}

func (om *ObstacleMap) addBaselineAround(centre Point) {
	left, right := om.baselineAround(centre)
	om.expandPath([]Point{left, right})
}

func (om *ObstacleMap) expandPath(points Path) {
	segments := zipWithNext(points)

	for _, segment := range segments {
		om.expandSegment(segment)
	}
}

func (om *ObstacleMap) isBlocked(point Point) bool {
	_, exists := om.points[point]
	return exists
}

func (om *ObstacleMap) expandSegment(segment Segment) {
	x1, x2 := inOrder(segment.start.x, segment.end.x)
	y1, y2 := inOrder(segment.start.y, segment.end.y)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			om.points[Point{x, y}] = ROCK
		}
	}
	if y2 > om.lowest {
		om.lowest = y2
	}
}

func (om *ObstacleMap) baselineAround(start Point) (Point, Point) {
	depth := om.lowest + 2
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
	points := strings.Split(input, " -> ")
	result := make([]Point, len(points))

	for i, coords := range points {
		x, y, _ := strings.Cut(coords, ",")
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
