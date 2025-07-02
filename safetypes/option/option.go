// Package option implements https://doc.rust-lang.org/std/option/enum.Option.html
package option

import (
	"reflect"
)

type Option[T any] interface {
	IsNone() bool
	IsNoneOr(func(T) bool) bool
	IsSome() bool
	IsSomeAnd(func(T) bool) bool
	Expect(msg string) T
	Unwrap() T
	UnwrapOr(def T) T
	UnwrapOrElse(f func() T) T
	UnwrapOrDefault() T
}

func From[T comparable](val T) Option[T] {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return none[T]{}
	}

	return Some(val)
}
