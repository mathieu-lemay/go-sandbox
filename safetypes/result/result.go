package result

import "fmt"

type Result[T any] struct {
	val T
	err error
}

func Ok[T any](val T) Result[T] {
	return Result[T]{
		val: val,
		err: nil,
	}
}

func Err(err error) Result[interface{}] {
	return Result[interface{}]{
		val: nil,
		err: err,
	}
}

func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(fmt.Errorf("%s: %w", msg, r.err))
	}

	return r.val
}

func (r Result[T]) ExpectErr(msg string) error {
	if r.err == nil {
		panic(fmt.Errorf("%s: %v", msg, r.val))
	}

	return r.err
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) IsErrAnd(predicate func(error) bool) bool {
	return r.err != nil && predicate(r.err)
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsOkAnd(predicate func(T) bool) bool {
	return r.err == nil && predicate(r.val)
}
