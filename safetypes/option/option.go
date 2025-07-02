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
	Inspect(f func(T)) Option[T]
}

func From[T comparable](val T) Option[T] {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return none[T]{}
	}

	return Some(val)
}

func Map[T any, U any](opt Option[T], f func(T) U) Option[U] {
	s, ok := opt.(some[T])
	if !ok {
		return none[U]{}
	}

	return Some(f(s.val))
}

func MapOr[T any, U any](opt Option[T], def U, f func(T) U) U {
	s, ok := opt.(some[T])
	if !ok {
		return def
	}

	return f(s.val)
}

func MapOrElse[T any, U any](opt Option[T], factory func() U, f func(T) U) U {
	s, ok := opt.(some[T])
	if !ok {
		return factory()
	}

	return f(s.val)
}
