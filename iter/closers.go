package betteriter

import (
	"strings"
)

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
