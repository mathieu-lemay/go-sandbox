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

func Map[T any, U any, E error](result Result[T, E], f func(T) U) Result[U, E] {
	return ok[U, E]{}
}

func MapOr[T any, U any, E error](result Result[T, E], def U, f func(T) U) Result[U, E] {
	return Ok(def)
}

func MapOrElse[T any, U any, E error](
	result Result[T, E],
	factory func() U,
	mapper func(T) U,
) Result[U, E] {
	return Ok(factory())
}

func MapErr[T any, E error, F error](result Result[T, E], f func(E) F) Result[T, F] {
	return errT[T, F]{}
}
