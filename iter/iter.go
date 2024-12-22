package betteriter

import (
	"errors"
	"iter"
)

type Iterator[T any] struct {
	it  iter.Seq2[T, error]
	err error
}

// CLOSERS

// TODO: make them a functions?
func (i *Iterator[T]) Collect() ([]T, error) {
	output := make([]T, 0)

	for v, err := range i.it {
		if err != nil {
			return nil, err
		}
		output = append(output, v)
	}

	return output, nil
}

type Tuple[T any, U any] struct {
	A T
	B U
}

func New[T any](values []T) *Iterator[*T] {
	return &Iterator[*T]{
		it: func(yield func(*T, error) bool) {
			for i := range values {
				if !yield(&values[i], nil) {
					return
				}
			}
		},
	}
}

func Repeat[T any](v T) *Iterator[T] {
	return &Iterator[T]{
		it: func(yield func(T, error) bool) {
			for {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}

func RepeatN[T any](v T, n int) *Iterator[T] {
	return &Iterator[T]{
		it: func(yield func(T, error) bool) {
			for i := 0; i < n; i++ {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}

func Incr() *Iterator[int] {
	return &Iterator[int]{
		it: func(yield func(int, error) bool) {
			for i := 0; ; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func IncrN(n int) *Iterator[int] {
	return &Iterator[int]{
		it: func(yield func(int, error) bool) {
			for i := 0; ; i += n {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func IncrFrom(start int) *Iterator[int] {
	return &Iterator[int]{
		it: func(yield func(int, error) bool) {
			for i := start; ; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func IncrNFrom(start int, n int) *Iterator[int] {
	return &Iterator[int]{
		it: func(yield func(int, error) bool) {
			for i := start; ; i += n {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func Range(start int, end int) *Iterator[int] {
	return &Iterator[int]{
		it: func(yield func(int, error) bool) {
			for i := start; i < end; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func Cycle[T any](values []T) *Iterator[*T] {
	return &Iterator[*T]{
		it: func(yield func(*T, error) bool) {
			for {
				for i := range values {
					if !yield(&values[i], nil) {
						return
					}
				}
			}
		},
	}
}

func Chain[T any](slices ...[]T) *Iterator[*T] {
	return &Iterator[*T]{
		it: func(yield func(*T, error) bool) {
			for _, s := range slices {
				for i := range s {
					if !yield(&s[i], nil) {
						return
					}
				}
			}
		},
	}
}

func Product[T any, U any](p []T, q []U) *Iterator[*Tuple[*T, *U]] {
	return &Iterator[*Tuple[*T, *U]]{
		it: func(yield func(*Tuple[*T, *U], error) bool) {
			for i := range p {
				for j := range q {
					t := Tuple[*T, *U]{&p[i], &q[j]}
					if !yield(&t, nil) {
						return
					}
				}
			}
		},
	}
}

func Zip[T any, U any](a []T, b []U) *Iterator[Tuple[T, U]] {
	return &Iterator[Tuple[T, U]]{
		it: func(yield func(Tuple[T, U], error) bool) {
			for idx := range a {
				if idx >= len(b) {
					return
				}

				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
	}
}

func ZipEq[T any, U any](a []T, b []U) *Iterator[Tuple[T, U]] {
	return &Iterator[Tuple[T, U]]{
		it: func(yield func(Tuple[T, U], error) bool) {
			if len(a) != len(b) {
				yield(Tuple[T, U]{}, errors.New("slices are not the same length"))
				return
			}

			for idx := range a {
				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
	}
}
