package betteriter

func Map[T any, U any](iterator *Iterator[T], f func(T) (U, error)) *Iterator[U] {
	it := func(yield func(U, error) bool) {
		for v, err := range iterator.it {
			if err != nil {
				var u U
				yield(u, err)
				return
			}

			if !yield(f(v)) {
				return
			}
		}
	}

	return &Iterator[U]{it, nil}
}
