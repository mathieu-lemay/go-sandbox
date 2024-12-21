package betteriter

type Enumerator[T any] struct {
	Idx   int
	Value T
}

func Enumerate[T any](iterator *Iterator[T]) *Iterator[Enumerator[T]] {
	it := func(yield func(Enumerator[T], error) bool) {
		idx := 0
		for v, err := range iterator.it {
			enum := Enumerator[T]{idx, v}

			if !yield(enum, err) || err != nil {
				return
			}

			idx += 1
		}
	}

	return &Iterator[Enumerator[T]]{it, nil}
}
