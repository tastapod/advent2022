package day14

type Point struct {
	x, y int
}

type Obstacle rune

const ROCK Obstacle = '#'

type ObstacleMap map[Point]Obstacle

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

func ExpandPath(points ...Point) ObstacleMap {
	segments := zipWithNext(points...)
	result := make(ObstacleMap)

	for _, segment := range segments {
		segment.expandToMap(result)
	}
	return result
}

func zipWithNext(points ...Point) []Segment {
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

// Try with channels

func (s *Segment) expandToChannel(points chan<- Point) {
	x1, x2 := inOrder(s.start.x, s.end.x)
	y1, y2 := inOrder(s.start.y, s.end.y)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			points <- Point{x, y}
		}
	}
}

func zipWithNextToChannel(points <-chan Point, segments chan<- Segment) {
	left := <-points

	for right := range points {
		segments <- Segment{left, right}
		left = right
	}
	close(segments)
}

// ExpandPathToChannel sets up a processing chain as follows:
//
// path of Points -> sequence of Segments (start-end points) -> expanded stream of Points
//
// These are collected into a map which is published once the expanded stream of Points ends.
func ExpandPathToChannel(path <-chan Point, results chan<- ObstacleMap) {
	segments := make(chan Segment)
	points := make(chan Point)

	go zipWithNextToChannel(path, segments)
	go expandSegments(segments, points)
	go collectPoints(points, results)
}

func expandSegments(segments chan Segment, points chan Point) {
	for segment := range segments {
		segment.expandToChannel(points)
	}
	close(points)
}

func collectPoints(points chan Point, results chan<- ObstacleMap) {
	result := make(ObstacleMap)

	for point := range points {
		result[point] = ROCK
	}

	results <- result
}
