package day14

import (
	"fmt"
	"github.com/tastapod/advent2022/pair"
	. "github.com/tastapod/advent2022/segment"
	"strconv"
	"strings"
)

type Path []Point

func (path Path) ToSegments() ([]Segment, error) {
	pairs := pair.ZipWithNext(path)
	result := make([]Segment, len(pairs))

	for i, points := range pairs {
		segment, err := NewSegment(points.First, points.Second)
		if err != nil {
			return nil, err
		}
		result[i] = segment
	}
	return result, nil
}

func ParsePaths(input []string) []Path {
	result := make([]Path, len(input))

	for i, line := range input {
		result[i] = ParsePath(line)
	}
	return result
}

func ParsePath(input string) Path {
	coords := strings.Split(input, " -> ")
	result := make([]Point, len(coords))

	for i, coord := range coords {
		x, y, _ := strings.Cut(coord, ",")
		result[i] = Point{X: toInt(x), Y: toInt(y)}
	}
	return result
}

func toInt(s string) int {
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}
	panic(fmt.Sprintf("'%s': cannot parse int", s))
}

type Obstacle rune

const (
	SAND Obstacle = 'o'
	ROCK Obstacle = '#'
)

type ObstacleMap map[Point]Obstacle

func (om ObstacleMap) PlotObstaclePath(path Path, obstacle Obstacle) error {
	segments, err := path.ToSegments()
	if err != nil {
		return err
	}

	for _, segment := range segments {
		for _, point := range segment.Points() {
			om[point] = obstacle
		}
	}
	return nil
}

type Cave struct {
	obstacles ObstacleMap
	deepest   int
}

var StartPoint = Point{X: 500}

func NewCaveFromStrings(paths []string) (Cave, error) {
	return NewCaveFromPaths(ParsePaths(paths))
}

func NewCaveFromPaths(paths []Path) (Cave, error) {
	result := NewCave()
	for _, path := range paths {
		err := result.ExpandPath(path)
		if err != nil {
			return Cave{}, err
		}
	}
	return result, nil
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
// or the start point and false if it falls through, past deepest point.
func (cave *Cave) DropSandFrom(start Point) (Point, bool) {
	sand := start
	if cave.isBlocked(sand) {
		return sand, false
	}

grains:
	for below := sand; sand.Y <= cave.deepest; {
		for _, xBelow := range []int{sand.X, sand.X - 1, sand.X + 1} {
			below = Point{X: xBelow, Y: sand.Y + 1}
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

func (cave *Cave) AddBaselineAround(centre Point) error {
	left, right := cave.BaselineAround(centre)

	if err := cave.ExpandPath([]Point{left, right}); err != nil {
		return err
	}
	return nil
}

func (cave *Cave) ExpandPath(path Path) error {
	if err := cave.obstacles.PlotObstaclePath(path, ROCK); err != nil {
		return err
	}

	cave.deepest = findDeepest(path, cave.deepest)
	return nil
}

func findDeepest(path Path, deepest int) int {
	for _, point := range path {
		if point.Y > deepest {
			deepest = point.Y
		}
	}
	return deepest
}

func (cave *Cave) isBlocked(point Point) bool {
	_, exists := cave.obstacles[point]
	return exists
}

func (cave *Cave) BaselineAround(start Point) (Point, Point) {
	depth := cave.deepest + 2
	return Point{X: start.X - depth - 1, Y: depth}, Point{X: start.X + depth + 1, Y: depth}
}
