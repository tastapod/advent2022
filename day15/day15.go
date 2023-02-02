package day15

import (
	"fmt"
	. "github.com/tastapod/advent2022/segment"
)

type Reading struct {
	Sensor, Beacon Point
}

func (r *Reading) Distance() int {
	return absDiff(r.Sensor.X, r.Beacon.X) + absDiff(r.Sensor.Y, r.Beacon.Y)
}

func (r *Reading) PointsOnRow(row int) Segment {
	yDist := absDiff(r.Sensor.Y, row)

	xDist := r.Distance() - yDist
	if xDist > 0 {
		segment, _ := NewSegment(
			Point{X: r.Sensor.X - xDist, Y: row},
			Point{X: r.Sensor.X + xDist, Y: row})
		return segment
	} else {
		return Segment{}
	}
}

func (r *Reading) IntersectsRow(row int) bool {
	return r.Sensor.Y-r.Distance() <= row && row <= r.Sensor.Y+r.Distance()
}

func absDiff(a, b int) int {
	if a < b {
		return b - a
	} else {
		return a - b
	}
}

func NewReading(sensorLine string) (Reading, error) {
	var xSensor, ySensor int
	var xBeacon, yBeacon int
	parts, err := fmt.Sscanf(sensorLine, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
		&xSensor, &ySensor,
		&xBeacon, &yBeacon)
	if err != nil || parts != 4 {
		return Reading{}, err
	}
	return Reading{Sensor: Point{X: xSensor, Y: ySensor}, Beacon: Point{X: xBeacon, Y: yBeacon}}, nil
}

func CountNonBeaconPoints(row int, readings []string) int {
	linePoints := make(map[Point]bool)
	beacons := make([]Point, len(readings))

	for _, readingLine := range readings {
		reading, err := NewReading(readingLine)
		if err != nil {
			return 0
		}
		if reading.IntersectsRow(row) {
			segment := reading.PointsOnRow(row)
			for _, point := range segment.Points() {
				linePoints[point] = true
			}
		}
		if reading.Beacon.Y == row {
			beacons = append(beacons, reading.Beacon)
		}
	}

	for _, beacon := range beacons {
		delete(linePoints, beacon)
	}

	return len(linePoints)
}
