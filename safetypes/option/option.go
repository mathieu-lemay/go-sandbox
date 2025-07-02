package option

import (
	"reflect"
)

type Option[T any] interface {
	IsNone() bool
	// IsNoneOr(func(T) bool) bool
	IsSome() bool
	// IsSomeAnd(func(T) bool) bool
}

var (
	_ Option[any] = some[any]{}
	_ Option[any] = none{}
)

type some[T any] struct {
	val T
}

type none struct{}

func From[T comparable](val T) Option[T] {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return None()
	}

	return Some(val)
}

func Some[T any](val T) Option[T] {
	return some[T]{
		val: val,
	}
}

func None() Option[any] {
	return none{}
}

func (s some[T]) IsNone() bool {
	return false
}

func (s some[T]) IsSome() bool {
	return true
}

func (n none) IsNone() bool {
	return true
}

func (n none) IsSome() bool {
	return false
}

// func (r Option[T]) Expect(msg string) T {
//     if r.err != nil {
//         panic(fmt.Errorf("%s: %w", msg, r.err))
//     }
//
//     return r.val
// }
//
// func (r Option[T]) ExpectErr(msg string) error {
//     if r.err == nil {
//         panic(fmt.Errorf("%s: %v", msg, r.val))
//     }
//
//     return r.err
// }
//
// func (r Option[T]) IsErr() bool {
//     return r.err != nil
// }
//
// func (r Option[T]) IsErrAnd(predicate func(error) bool) bool {
//     return r.err != nil && predicate(r.err)
// }
//
// func (r Option[T]) IsOk() bool {
//     return r.err == nil
// }
//
// func (r Option[T]) IsOkAnd(predicate func(T) bool) bool {
//     return r.err == nil && predicate(r.val)
// }
