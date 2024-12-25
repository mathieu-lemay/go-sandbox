package betteriter

import (
	"strings"
)

func Collect[T any](iter *Iterator[T]) ([]T, error) {
	var output []T

	for v, err := range iter.it {
		if err != nil {
			return nil, err
		}
		output = append(output, v)
	}

	return output, nil
}

func Copied[T any](iter *Iterator[*T]) ([]T, error) {
	var output []T

	for v, err := range iter.it {
		if err != nil {
			return nil, err
		}
		output = append(output, *v)
	}

	return output, nil
}

func Any[T any](iter *Iterator[T]) (bool, error) {
	for _, err := range iter.it {
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func Fold[T any, U any](iter *Iterator[T], init U, adder func(cur U, item T) U) (U, error) {
	current := init

	for v, err := range iter.it {
		if err != nil {
			return current, err
		}

		current = adder(current, v)
	}

	return current, nil
}

func StringJoin(iter *Iterator[*string], s string) (string, error) {
	first := true
	sb := strings.Builder{}

	for v, err := range iter.it {
		if err != nil {
			return "", err
		}

		if !first {
			sb.WriteString(s)
		} else {
			first = false
		}

		sb.WriteString(*v)
	}

	return sb.String(), nil
}
