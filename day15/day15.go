package day15

import (
	"fmt"
	. "github.com/tastapod/advent2022/segment"
	"sort"
)

type Reading struct {
	Sensor, Beacon Point
}

func (r *Reading) Distance() int {
	return absDiff(r.Sensor.X, r.Beacon.X) + absDiff(r.Sensor.Y, r.Beacon.Y)
}

func (r *Reading) SegmentOnRow(row int) Segment {
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

// KleePoint used in Klee's algorithm below
type KleePoint struct {
	X     int
	IsEnd bool
}

// CountNonBeaconPoints uses [Klee's algorithm] as described at [Open Genus].
//
// [Klee's algorithm]: https://en.wikipedia.org/wiki/Klee%27s_measure_problem
// [Open Genus]: https://iq.opengenus.org/klee-algorithm/
func CountNonBeaconPoints(row int, readings []string) int {
	kleePoints := make([]KleePoint, 0)

	// build the Klee vector
	for _, readingLine := range readings {
		reading, err := NewReading(readingLine)
		if err != nil {
			return 0
		}
		if reading.IntersectsRow(row) {
			segment := reading.SegmentOnRow(row)
			kleePoints = append(kleePoints,
				KleePoint{segment.Start.X, false},
				KleePoint{segment.End.X, true})
		}
	}

	// sort it
	sort.Slice(kleePoints, func(i, j int) bool {
		left, right := kleePoints[i], kleePoints[j]
		if left.X < right.X {
			return true
		} else if left.X > right.X {
			return false
		} else {
			return left.IsEnd
		}
	})

	// calculate the size
	result := 0
	segmentDepth := 1 // how many segments deep at this point

	for i := 1; i < len(kleePoints); i++ {
		prev, this := kleePoints[i-1], kleePoints[i]
		if diff := this.X - prev.X; segmentDepth > 0 && diff > 0 {
			result += diff
		}
		if this.IsEnd {
			segmentDepth -= 1
		} else {
			segmentDepth += 1
		}
	}
	return result
}
