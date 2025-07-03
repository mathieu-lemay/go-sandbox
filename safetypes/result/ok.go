// Package result implements https://doc.rust-lang.org/std/result/enum.Result.html
package result

import "fmt"

func Ok[T any](val T) Result[T, error] {
	return ok[T, error]{
		val: val,
	}
}

type ok[T any, E error] struct {
	val T
}

func (o ok[T, E]) IsOk() bool {
	return true
}

func (o ok[T, E]) IsOkAnd(f func(T) bool) bool {
	return f(o.val)
}

func (o ok[T, E]) IsErr() bool {
	return false
}

func (o ok[T, E]) IsErrAnd(_ func(error) bool) bool {
	return false
}

func (o ok[T, E]) Expect(_ string) T {
	return o.val
}

func (o ok[T, E]) ExpectErr(msg string) E {
	panic(fmt.Errorf("%s: %v", msg, o.val))
}

func (o ok[T, E]) Unwrap() T {
	return o.val
}

func (o ok[T, E]) UnwrapOr(_ T) T {
	return o.val
}

func (o ok[T, E]) UnwrapOrElse(_ func() T) T {
	return o.val
}

func (o ok[T, E]) UnwrapOrDefault() T {
	return o.val
}

func (o ok[T, E]) UnwrapErr() E {
	panic(fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", o.val))
}
