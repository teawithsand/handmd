package single

import "context"

// Wraps single value of type T, which may be loaded on demand.
// It's something like Lazy in languages like kotlin.
//
// It's similar to iterable, but with single item.
type Single[T any] interface {
	Obtain(ctx context.Context) (res T, err error)
}

type SingleFunc[T any] func(ctx context.Context) (res T, err error)

func (f SingleFunc[T]) Obtain(ctx context.Context) (res T, err error) {
	return f(ctx)
}

// Memorizes single provided in memory.
// Once value is loaded, underling single's Obtain isn't called anymore.
// Errors are not memorized.
func Memorize[T any](s Single[T]) Single[T] {
	var memorizedValue T
	isLoaded := false
	return SingleFunc[T](func(ctx context.Context) (res T, err error) {
		if isLoaded {
			res = memorizedValue
			return
		}

		res, err = s.Obtain(ctx)
		if err != nil {
			return
		}

		memorizedValue = res
		isLoaded = true
		return
	})
}

func Eager[T any](data T) Single[T] {
	return SingleFunc[T](func(ctx context.Context) (res T, err error) {
		res = data
		return
	})
}

func Map[T any, E any](s Single[T], mapper func(ctx context.Context, value T) (res E, err error)) (res Single[E]) {
	res = SingleFunc[E](func(ctx context.Context) (res E, err error) {
		imm, err := s.Obtain(ctx)
		if err != nil {
			return
		}

		res, err = mapper(ctx, imm)
		if err != nil {
			return
		}

		return
	})
	return
}
