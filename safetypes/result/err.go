// Package result implements https://doc.rust-lang.org/std/result/enum.Result.html
package result

func Err[E error](err E) Result[any, E] {
	return errT[any, E]{
		err: err,
	}
}

type errT[T any, E error] struct {
	err E
}

func (e errT[T, E]) IsOk() bool {
	return false
}

func (e errT[T, E]) IsOkAnd(_ func(T) bool) bool {
	return false
}

func (e errT[T, E]) IsErr() bool {
	return true
}

func (e errT[T, E]) IsErrAnd(f func(error) bool) bool {
	return f(e.err)
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
