// Package safetypes implements Rust-style Option[T] and Result[T, E]
package safetypes

import (
	"github.com/mathieu-lemay/go-sandbox/safetypes/option"
	"github.com/mathieu-lemay/go-sandbox/safetypes/result"
)

// AsOkOr converts an option to a result.Ok when opt is option.Some or result.Err when opt is option.None.
func AsOkOr[T any](opt option.Option[T], err error) result.Result[T, error] {
	if opt.IsSome() {
		return result.Ok(opt.Unwrap())
	}

	var v T

	return result.From(v, err)
}

// AsOkOrElse converts an option to a result.Ok when opt is option.Some or result.Err when opt is option.None.
func AsOkOrElse[T any](opt option.Option[T], f func() error) result.Result[T, error] {
	if opt.IsSome() {
		return result.Ok(opt.Unwrap())
	}

	var v T

	return result.From(v, f())
}

// AsOptionValue converts a result.Result to an option.Some when res is result.Ok or option.None when res is result.Err.
func AsOptionValue[T any, E error](res result.Result[T, E]) option.Option[T] {
	if res.IsOk() {
		return option.Some(res.Unwrap())
	}

	var v T

	return option.From(v)
}

// AsOptionErr converts a result.Result to an option.Some when res is result.Err or option.None when res is result.Ok.
func AsOptionErr[T any, E error](res result.Result[T, E]) option.Option[E] {
	if res.IsErr() {
		return option.From(res.UnwrapErr())
	}

	var v E

	return option.From(v)
}
