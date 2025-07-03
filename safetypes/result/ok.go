// Package result implements https://doc.rust-lang.org/std/result/enum.Result.html
package result

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

//func (r Result[T, E]) Expect(msg string) T {
//	if r.err != nil {
//		panic(fmt.Errorf("%s: %w", msg, r.err))
//	}
//
//	return r.val
//}
//
//func (r Result[T, E]) ExpectErr(msg string) error {
//	if r.err == nil {
//		panic(fmt.Errorf("%s: %v", msg, r.val))
//	}
//
//	return r.err
//}
//
//func (r Result[T, E]) IsErr() bool {
//	return r.err != nil
//}
//
//func (r Result[T, E]) IsErrAnd(predicate func(error) bool) bool {
//	return r.err != nil && predicate(r.err)
//}
//
//func (r Result[T, E]) IsOk() bool {
//	return r.err == nil
//}
//
//func (r Result[T, E]) IsOkAnd(predicate func(T) bool) bool {
//	return r.err == nil && predicate(r.val)
//}
