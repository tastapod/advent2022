package pair

type Pair[T any] struct {
	First, Second T
}

func ZipWithNext[T any](ts []T) []Pair[T] {
	ends := ts[1:]
	result := make([]Pair[T], len(ends))

	for i, end := range ends {
		result[i] = Pair[T]{First: ts[i], Second: end}
	}
	return result
}
