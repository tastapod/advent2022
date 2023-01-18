package day14

import (
	"reflect"
	"strconv"
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

	pairs := zipWithNext(Point{10, 20}, Point{25, 35}, Point{40, 50}, Point{55, 65}, Point{70, 80})

	if !reflect.DeepEqual(pairs, expected) {
		errorNotEqual(t, expected, pairs)
	}
}

func errorNotEqual[T any](t *testing.T, expected T, actual T) {
	t.Error("Expected:\n", expected, "\nbut got\n", actual)
}

func TestExpandsPath(t *testing.T) {
	points := ExpandPath(Point{498, 4}, Point{498, 6}, Point{496, 6})

	for _, point := range []Point{{498, 4}, {498, 5}, {498, 6}, {497, 6}, {496, 6}} {
		_, found := points[point]
		if !found {
			t.Error(point, " expected but not found")
		}
	}

	for _, point := range []Point{{497, 4}, {500, 5}, {498, 10}} {
		_, found := points[point]
		if found {
			t.Error(point, " found but not expected")
		}
	}
}

func TestExpandsPathWithChannels(t *testing.T) {
	path := make(chan Point)
	obstacleMaps := make(chan ObstacleMap)

	go ExpandPathToChannel(path, obstacleMaps)

	for _, point := range []Point{{498, 4}, {498, 6}, {496, 6}} {
		path <- point
	}
	close(path)

	obstacles := <-obstacleMaps

	for _, point := range []Point{{498, 4}, {498, 5}, {498, 6}, {497, 6}, {496, 6}} {
		_, found := obstacles[point]
		if !found {
			t.Error(point, " expected but not found")
		}
	}

	for _, point := range []Point{{497, 4}, {500, 5}, {498, 10}} {
		_, found := obstacles[point]
		if found {
			t.Error(point, " found but not expected")
		}
	}
}

func TestParsesPath(t *testing.T) {
	input := "503,4 -> 502,4 -> 502,9 -> 494,9"
	expected := []Point{{503, 4}, {502, 4}, {502, 9}, {494, 9}}

	if points := ParseLine(input); !reflect.DeepEqual(expected, points) {
		errorNotEqual(t, expected, points)
	}
}

func ParseLine(input string) []Point {
	points := strings.Split(input, " -> ")
	result := make([]Point, len(points))

	for i, coords := range points {
		x, y, _ := strings.Cut(coords, ",")
		result[i] = Point{toInt(x), toInt(y)}
	}
	return result
}

func toInt(s string) int {
	if result, err := strconv.Atoi(s); err == nil {
		return result
	}
	panic(s + ": cannot parse int")
}
