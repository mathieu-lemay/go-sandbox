package betteriter

func (i *Iterator[T, U]) Map(f func(T) (U, error)) *Iterator[*U, U] {
	it := func(yield func(*U, error) bool) {
		for v, err := range i.it {
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

	return &Iterator[*U, U]{it}
}
