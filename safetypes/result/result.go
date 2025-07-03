// Package result implements https://doc.rust-lang.org/std/result/enum.Result.html
package result

import (
	"reflect"
)

type Result[T any, E error] interface {
	IsOk() bool
	IsOkAnd(f func(T) bool) bool
	IsErr() bool
	IsErrAnd(f func(error) bool) bool
	//Ok() option.Option[T]
	//Err() option.Option[E]
	//Inspect(func(*T)) Result[T, E]
	//InspectErr(func(*E)) Result[T, E]
	Expect(msg string) T
	ExpectErr(msg string) E
	Unwrap() T
	UnwrapOr(def T) T
	UnwrapOrElse(f func() T) T
	UnwrapOrDefault() T
	UnwrapErr() E
}

func From[T any, E error](val T, err E) Result[T, E] {
	if !reflect.ValueOf(&err).Elem().IsNil() {
		return errT[T, E]{
			err: err,
		}
	}

	return ok[T, E]{
		val: val,
	}
}

func Map[T any, U any, E error](res Result[T, E], f func(T) U) Result[U, E] {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return errT[U, E]{res.UnwrapErr()}
	}

	val := f(s.val)

	return ok[U, E]{val}
}

func MapOr[T any, U any, E error](res Result[T, E], def U, f func(T) U) U {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return def
	}

	return f(s.val)
}

func MapOrElse[T any, U any, E error](
	res Result[T, E],
	factory func() U,
	mapper func(T) U,
) U {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return factory()
	}

	return mapper(s.val)
}

func MapErr[T any, E error, F error](res Result[T, E], f func(E) F) Result[T, F] {
	s, isOk := res.(ok[T, E])
	if isOk {
		return ok[T, F]{s.val}
	}

	return errT[T, F]{f(res.UnwrapErr())}
}
