package option

import (
	"errors"

	"github.com/mathieu-lemay/go-sandbox/safetypes/result"
)

var (
	_ Option[any] = none[any]{}
)

type none[T any] struct{}

func None() Option[any] {
	return none[any]{}
}

func (n none[T]) IsNone() bool {
	return true
}

func (n none[T]) IsNoneOr(func(T) bool) bool {
	return true
}

func (n none[T]) IsSome() bool {
	return false
}

func (n none[T]) IsSomeAnd(func(T) bool) bool {
	return false
}

func (n none[T]) Expect(msg string) T {
	panic(errors.New(msg))
}

func (n none[T]) Unwrap() T {
	panic(errors.New("called `Option.Unwrap()` on a `None` value"))
}

func (n none[T]) UnwrapOr(def T) T {
	return def
}

func (n none[T]) UnwrapOrElse(f func() T) T {
	return f()
}

func (n none[T]) UnwrapOrDefault() T {
	var v T
	return v
}

func (n none[T]) Inspect(_ func(T)) Option[T] {
	return n
}

func (n none[T]) OkOr(err error) result.Result[T] {
	var v T
	return result.From(v, err)
}

func (n none[T]) OkOrElse(f func() error) result.Result[T] {
	var v T
	return result.From(v, f())
}

func (n none[T]) Filter(_ func(T) bool) Option[T] {
	return n
}

func (n none[T]) Or(other Option[T]) Option[T] {
	return other
}

func (n none[T]) OrElse(f func() Option[T]) Option[T] {
	return f()
}

func (n none[T]) Xor(other Option[T]) Option[T] {
	return other
}
