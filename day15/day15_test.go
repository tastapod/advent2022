package day15_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/tastapod/advent2022/day15"
	"strings"
	"testing"
)

func TestParsesSensorLine(t *testing.T) {
	assert := assert.New(t)

	reading, err := NewReading("Sensor at x=2, y=18: closest beacon is at x=-2, y=15")
	assert.Equal(Point{X: 2, Y: 18}, reading.Sensor)
	assert.Equal(Point{X: -2, Y: 15}, reading.Beacon)
	assert.Equal(nil, err)

	reading, err = NewReading("Sensor at some other place")
	assert.Equal(Reading{}, reading)
	assert.NotEqual(nil, err)
}

func reading(sensorX, sensorY, beaconX, beaconY int) Reading {
	return Reading{
		Sensor: Point{X: sensorX, Y: sensorY},
		Beacon: Point{X: beaconX, Y: beaconY},
	}
}

func TestCalculatesManhattanDistance(t *testing.T) {
	reading := reading(2, 18, -2, 15)
	assert.Equal(t, 7, reading.Radius())
}

func TestFindsPointsOnGivenRow(t *testing.T) {
	assert := assert.New(t)

	reading := reading(8, 7, 2, 10)
	span, intersects := reading.SpanOnRow(10)

	assert.True(intersects)
	assert.Equal(Span{Start: 2, End: 14}, span)
}

var sampleReadingLines = strings.Split(strings.TrimSpace(`
Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3
`), "\n")

func parseReadings(t *testing.T) []Reading {
	readings, err := ParseReadings(sampleReadingLines)
	if err != nil {
		t.Error(err)
	}
	return readings
}

func TestCountsPointsIntersectingRow(t *testing.T) {
	readings := parseReadings(t)
	assert.Equal(t, 26, SumOverlappingSegmentLengths(10, readings))
}

func TestFindsPossiblePoints(t *testing.T) {
	readings := parseReadings(t)
	point, _ := FindVacantPoint(20, readings)
	assert.Equal(t, Point{X: 14, Y: 11}, point)
}

func TestCalculatesTuningFrequency(t *testing.T) {
	assert.Equal(t, 56000011, TuningFrequency(Point{X: 14, Y: 11}))
}
