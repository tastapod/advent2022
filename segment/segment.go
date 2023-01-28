package segment

import (
	"errors"
	"fmt"
)

type Point struct {
	X, Y int
}

type Segment struct {
	Start, End Point
	points     []Point
}

// NewSegment creates a Segment from two Points in a straight line.
//
// A horizontal segment has Start and End from left to right, and a vertical segment from top to bottom.
// It returns an error if the start and end points are not horizontally or vertically in line with each other.
func NewSegment(start Point, end Point) (Segment, error) {
	if start.X != end.X && start.Y != end.Y {
		return Segment{}, errors.New(fmt.Sprintf("%v -> %v is not a straight line", start, end))
	}
	x1, x2 := inOrder(start.X, end.X)
	y1, y2 := inOrder(start.Y, end.Y)
	points := make([]Point, 0)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			points = append(points, Point{x, y})
		}
	}
	return Segment{Start: Point{x1, y1}, End: Point{X: x2, Y: y2}, points: points}, nil
}

func (s *Segment) Points() []Point {
	result := make([]Point, len(s.points))
	copy(result, s.points)
	return result
}

func (s *Segment) Contains(point Point) bool {
	for _, maybe := range s.points {
		if point == maybe {
			return true
		}
	}
	return false
}

func inOrder(i1 int, i2 int) (int, int) {
	if i1 <= i2 {
		return i1, i2
	} else {
		return i2, i1
	}
}
