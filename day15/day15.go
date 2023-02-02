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

func ParseReadings(readingLines []string) ([]Reading, error) {
	readings := make([]Reading, 0, len(readingLines))
	for _, readingLine := range readingLines {
		reading, err := NewReading(readingLine)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, nil
}

// SumOverlappingSegmentLengths uses [Klee's algorithm] as described at [Open Genus].
//
// [Klee's algorithm]: https://en.wikipedia.org/wiki/Klee%27s_measure_problem
// [Open Genus]: https://iq.opengenus.org/klee-algorithm/
func SumOverlappingSegmentLengths(row int, readings []Reading) int {
	kleePoints := buildKleeVector(row, readings)
	return calculateOverlappingLength(kleePoints)
}

func FindVacantPoint(limit int, readings []Reading) *Point {
	for y := 0; y <= limit; y++ {
		kleeVector := buildKleeVector(y, readings)
		if x := findEmptyPoint(kleeVector); x != nil {
			result := Point{X: *x, Y: y}
			return &result
		}
	}
	return nil
}

func TuningFrequency(point Point) int {
	return 4_000_000*point.X + point.Y
}

// findEmptyPoint is a riff on Klee's algorithm where we look for any
// non-zero run where segment depth is 0, i.e. not covered by a segment
//
// Return an int pointer so we can use nil for not found
func findEmptyPoint(kleeVector []KleePoint) *int {
	segmentDepth := 1 // how many segments deep at this point

	for i := 1; i < len(kleeVector); i++ {
		prev, this := kleeVector[i-1], kleeVector[i]
		if diff := this.X - prev.X; segmentDepth == 0 && diff > 0 {
			// we found one!
			result := prev.X + 1
			return &result
		}
		if this.IsEnd {
			segmentDepth -= 1
		} else {
			segmentDepth += 1
		}
	}
	return nil
}

func buildKleeVector(row int, readings []Reading) []KleePoint {
	result := make([]KleePoint, 0)

	for _, reading := range readings {
		if reading.IntersectsRow(row) {
			segment := reading.SegmentOnRow(row)
			result = append(result,
				KleePoint{segment.Start.X, false},
				KleePoint{segment.End.X, true})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		if left.X < right.X {
			return true
		} else if left.X > right.X {
			return false
		} else {
			return left.IsEnd
		}
	})
	return result
}

func calculateOverlappingLength(kleeVector []KleePoint) int {
	result := 0
	segmentDepth := 1 // how many segments deep at this point

	for i := 1; i < len(kleeVector); i++ {
		prev, this := kleeVector[i-1], kleeVector[i]
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
