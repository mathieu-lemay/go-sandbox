package betteriter

func (i *Iterator[T, U]) Collect() ([]T, error) {
	output := make([]T, 0)

	for v, err := range i.it {
		if err != nil {
			return nil, err
		}
		output = append(output, v)
	}

	return output, nil
}

func (i *Iterator[T, U]) Any(predicate func(T) bool) (bool, error) {
	for v, err := range i.it {
		if err != nil {
			return false, err
		}

		if predicate(v) {
			return true, nil
		}
	}

	return false, nil
}

func (i *Iterator[T, U]) All(predicate func(T) bool) (bool, error) {
	for v, err := range i.it {
		if err != nil {
			return false, err
		}

		if !predicate(v) {
			return false, nil
		}
	}

	return true, nil
}

func (i *Iterator[T, U]) Fold(init U, adder func(cur U, item T) U) (U, error) {
	current := init

	for v, err := range i.it {
		if err != nil {
			return current, err
		}

		current = adder(current, v)
	}

	return current, nil
}
