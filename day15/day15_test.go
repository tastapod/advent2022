package day15_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/tastapod/advent2022/day15"
	. "github.com/tastapod/advent2022/segment"
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

func TestCalculatesManhattanDistance(t *testing.T) {
	reading := Reading{Sensor: Point{X: 2, Y: 18}, Beacon: Point{X: -2, Y: 15}}
	assert.Equal(t, 7, reading.Distance())
}

func TestFindsPointsOnGivenRow(t *testing.T) {
	assert := assert.New(t)

	reading := Reading{Sensor: Point{X: 8, Y: 7}, Beacon: Point{X: 2, Y: 10}}
	segment := reading.SegmentOnRow(10)

	assert.Equal(13, len(segment.Points()))
	assert.True(segment.Contains(Point{X: 2, Y: 10}))
	assert.False(segment.Contains(Point{X: 1, Y: 10}))
}

func TestChecksForIntersectingRow(t *testing.T) {
	assert := assert.New(t)

	reading := Reading{Sensor: Point{X: 8, Y: 7}, Beacon: Point{X: 2, Y: 10}}
	assert.False(reading.IntersectsRow(-3))
	assert.True(reading.IntersectsRow(-2))
	assert.True(reading.IntersectsRow(16))
	assert.False(reading.IntersectsRow(17))
}

var sampleReadings = strings.Split(strings.TrimSpace(`
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

func TestCountsPointsIntersectingRow(t *testing.T) {
	assert.Equal(t, 26, CountNonBeaconPoints(10, sampleReadings))
}
