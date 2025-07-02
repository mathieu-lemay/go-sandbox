package option

var (
	_ Option[any] = some[any]{}
)

type some[T any] struct {
	val T
}

func Some[T any](val T) Option[T] {
	return some[T]{
		val: val,
	}
}

func (s some[T]) IsNone() bool {
	return false
}

func (s some[T]) IsNoneOr(f func(T) bool) bool {
	return f(s.val)
}

func (s some[T]) IsSome() bool {
	return true
}

func (s some[T]) IsSomeAnd(f func(T) bool) bool {
	return f(s.val)
}

func (s some[T]) Expect(_ string) T {
	return s.val
}

func (s some[T]) Unwrap() T {
	return s.val
}

func (s some[T]) UnwrapOr(_ T) T {
	return s.val
}

func (s some[T]) UnwrapOrElse(_ func() T) T {
	return s.val
}

func (s some[T]) UnwrapOrDefault() T {
	return s.val
}

func (s some[T]) Inspect(f func(T)) Option[T] {
	f(s.val)
	return s
}
