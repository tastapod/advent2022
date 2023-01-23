package day14

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
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
	cave := Cave{Point{500, 0}, NewObstacleMap()}
	for _, path := range ParsePaths(paths) {
		cave.obstacles.ExpandPath(path)
	}
	return cave
}

// dropSand drops a grain of sand from the start point until it settles.
// It returns the final resting place of the sand and true if the sand settled,
// or the start point and false if it fell through.
func (c *Cave) dropSand() (Point, bool) {
	bottom := c.bottom()
	sand := c.start

	if c.obstacles.isBlocked(sand) {
		return sand, false
	}

grains:
	for sand.y <= bottom {
		for _, xBelow := range []int{sand.x, sand.x - 1, sand.x + 1} {
			below := Point{xBelow, sand.y + 1}
			if !c.obstacles.isBlocked(below) {
				sand = below
				continue grains
			}
		}
		// we have come to rest
		c.obstacles[sand] = SAND
		return sand, true
	}
	// we fell through
	return c.start, false
}

func (c *Cave) bottom() int {
	result := 0
	for rock := range c.obstacles {
		if rock.y > result {
			result = rock.y
		}
	}
	return result
}

func (c *Cave) baseline() (Point, Point) {
	bottom := c.bottom() + 2
	return Point{c.start.x - bottom - 1, bottom}, Point{c.start.x + bottom + 1, bottom}
}

func (c *Cave) FillWithSand() int {
	before := len(c.obstacles)

	for _, settled := c.dropSand(); settled; _, settled = c.dropSand() {

	}
	after := len(c.obstacles)
	return after - before
}

func (c *Cave) AddBaseline() {
	start, end := c.baseline()
	c.obstacles.ExpandPath([]Point{start, end})
}

type ObstacleMap map[Point]Obstacle

func NewObstacleMap() ObstacleMap {
	return make(ObstacleMap)
}

func (om ObstacleMap) ExpandPath(points Path) {
	segments := zipWithNext(points)

	for _, segment := range segments {
		segment.expandToMap(om)
	}
}

func (om ObstacleMap) isBlocked(point Point) bool {
	_, exists := om[point]
	return exists
}

type Segment struct {
	start, end Point
}

func (s *Segment) expandToMap(obstacleMap ObstacleMap) {
	x1, x2 := inOrder(s.start.x, s.end.x)
	y1, y2 := inOrder(s.start.y, s.end.y)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			obstacleMap[Point{x, y}] = ROCK
		}
	}
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

func ParsePath(input string) Path {
	points := strings.Split(input, " -> ")
	result := make([]Point, len(points))

	for i, coords := range points {
		x, y, _ := strings.Cut(coords, ",")
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
