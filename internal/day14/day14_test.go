package day14

import (
	"reflect"
	"testing"
)

func TestCreatesPairs(t *testing.T) {
	var expected = []Segment{
		{Point{10, 20}, Point{25, 35}},
		{Point{25, 35}, Point{40, 50}},
		{Point{40, 50}, Point{55, 65}},
		{Point{55, 65}, Point{70, 80}},
	}

	var pairs = zipWithNext(Point{10, 20}, Point{25, 35}, Point{40, 50}, Point{55, 65}, Point{70, 80})

	if !reflect.DeepEqual(pairs, expected) {
		t.Error("Expected:\n", expected, "\nbut got\n", pairs)
	}
}

func TestFollowsPath(t *testing.T) {
	var points = ExpandPath(Point{498, 4}, Point{498, 6}, Point{496, 6})

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

func TestFollowsPathWithChannels(t *testing.T) {
	path := make(chan Point)
	results := make(chan ObstacleMap)

	go ExpandPathToChannel(path, results)

	for _, point := range []Point{{498, 4}, {498, 6}, {496, 6}} {
		path <- point
	}
	close(path)

	obstacles := <-results

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
