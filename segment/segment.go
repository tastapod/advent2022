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
	return Segment{Start: Point{x1, y1}, End: Point{X: x2, Y: y2}}, nil
}

func (s *Segment) Points() []Point {
	result := make([]Point, 0)
	for x := s.Start.X; x <= s.End.X; x++ {
		for y := s.Start.Y; y <= s.End.Y; y++ {
			result = append(result, Point{x, y})
		}
	}
	return result
}

func (s *Segment) Contains(point Point) bool {
	return s.Start.X <= point.X && point.X <= s.End.X &&
		s.Start.Y <= point.Y && point.Y <= s.End.Y
}

func inOrder(i1 int, i2 int) (int, int) {
	if i1 <= i2 {
		return i1, i2
	} else {
		return i2, i1
	}
}
