package day15

import (
	"fmt"
	"runtime"
	"sort"
)

type Reading struct {
	Sensor, Beacon Point
}

type Point struct {
	X, Y int
}

type Span struct {
	Start, End int
}

func (p *Point) ManhattanDist(other Point) int {
	return absDiff(p.X, other.X) + absDiff(p.Y, other.Y)
}

func (r *Reading) Radius() int {
	return r.Sensor.ManhattanDist(r.Beacon)
}

// SpanOnRow finds the Span of the X coordinates where a Sensor intersects a row.
// If the Sensor does not intersect the row, return zero Span and false.
func (r *Reading) SpanOnRow(row int) (Span, bool) {
	yDist := absDiff(r.Sensor.Y, row)

	if xDist := r.Radius() - yDist; xDist >= 0 {
		return Span{r.Sensor.X - xDist, r.Sensor.X + xDist}, true
	} else {
		return Span{}, false
	}
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

func FindVacantPoint(limit int, readings []Reading) (Point, bool) {
	return findVacantPointsWithGoroutines(limit, readings), true
}

func findVacantPointsSerially(limit int, readings []Reading) (Point, bool) {
	beaconsByRow := mapOfBeaconsByRow(readings)

	// for each row
	for row := 0; row < limit; row++ {
		spans := make([]Span, 0)

		if beaconSpans, found := beaconsByRow[row]; found {
			copy(spans, beaconSpans)
		}

		for _, reading := range readings {
			if span, intersectsRow := reading.SpanOnRow(row); intersectsRow {
				spans = append(spans, span)
			}
		}

		// sort ranges by start point
		sort.Slice(spans, func(i, j int) bool { return spans[i].Start < spans[j].Start })

		// for each span 1-(n-1)
		spanSoFar := spans[0]

		for _, span := range spans[1:] {

			// if span overlaps span so far
			if span.Start <= spanSoFar.End {
				if span.End > spanSoFar.End {
					spanSoFar = Span{spanSoFar.Start, span.End}
				}
			} else {
				// FOUND EMPTY SPOT!
				return Point{spanSoFar.End + 1, row}, true
			}
		}
	}

	// didn't find the vacant point
	return Point{}, false
}

// Will block forever if there is no answer
func findVacantPointsWithGoroutines(limit int, readings []Reading) Point {
	beaconsByRow := mapOfBeaconsByRow(readings)

	results := make(chan Point)

	numWorkers := runtime.NumCPU()

	for worker := 0; worker < numWorkers; worker++ {
		go func(startRow int, results chan<- Point) {
			for row := startRow; row < limit; row += numWorkers {
				spans := make([]Span, 0)

				if beaconSpans, found := beaconsByRow[row]; found {
					copy(spans, beaconSpans)
				}

				for _, reading := range readings {
					if span, intersectsRow := reading.SpanOnRow(row); intersectsRow {
						spans = append(spans, span)
					}
				}

				// sort ranges by start point
				sort.Slice(spans, func(i, j int) bool { return spans[i].Start < spans[j].Start })

				// for each span 1-(n-1)
				spanSoFar := spans[0]

				for _, span := range spans[1:] {

					// if span overlaps span so far
					if span.Start <= spanSoFar.End {
						if span.End > spanSoFar.End {
							spanSoFar = Span{spanSoFar.Start, span.End}
						}
					} else {
						// FOUND EMPTY SPOT!
						results <- Point{spanSoFar.End + 1, row}
					}
				}
			}
		}(worker, results)
	}

	return <-results
}

func mapOfBeaconsByRow(readings []Reading) map[int][]Span {
	result := make(map[int][]Span)
	for _, reading := range readings {
		row := reading.Beacon.Y
		span := Span{reading.Beacon.X, reading.Beacon.X}

		if beacons, found := result[row]; found {
			result[row] = append(beacons, span)
		} else {
			result[row] = []Span{span}
		}
	}
	return result
}

func TuningFrequency(point Point) int {
	return 4_000_000*point.X + point.Y
}

func buildKleeVector(row int, readings []Reading) []KleePoint {
	result := make([]KleePoint, 0)

	for _, reading := range readings {
		if span, intersects := reading.SpanOnRow(row); intersects {
			result = append(result,
				KleePoint{span.Start, false},
				KleePoint{span.End, true})
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
