package betteriter

func Map[T any, U any](iterator *Iterator[T], f func(T) (U, error)) *Iterator[*U] {
	it := func(yield func(*U, error) bool) {
		for v, err := range iterator.it {
			if err != nil {
				yield(nil, err)
				return
			}

			u, err := f(v)
			if !yield(&u, err) {
				return
			}
		}
	}

	return &Iterator[*U]{it}
}
