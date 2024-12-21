package betteriter

func FilterMap[T any, U any](iterator *Iterator[T], f func(T) (U, bool, error)) *Iterator[U] {
	inner := func(yield func(U, error) bool) {
		for v, err := range iterator.it {
			if err != nil {
				var u U
				yield(u, err)
				return
			}

			u, ok, err := f(v)
			if ok || err != nil {
				if !yield(u, err) {
					return
				}
			}
		}
	}

	return &Iterator[U]{inner, nil}
}
