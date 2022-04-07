package iter

import "context"

// Collection of elements of type T, which may be iterated on.
// Note: iterables are allowed to be concurrent. They may launch multiple goroutines at the same time, even though that's not the default.
type Iterable[T any] interface {
	Iterate(ctx context.Context, recv Receiver[T]) (err error)
}

type IteratorFunc[T any] func(ctx context.Context, recv Receiver[T]) (err error)

func (f IteratorFunc[T]) Iterate(ctx context.Context, recv Receiver[T]) (err error) {
	return f(ctx, recv)
}

func Collect[T any](ctx context.Context, it Iterable[T]) (res []T, err error) {
	err = it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
		res = append(res, data)
		return
	}))
	return
}

func Slice[T any](data []T) Iterable[T] {
	return IteratorFunc[T](func(ctx context.Context, recv Receiver[T]) (err error) {
		for _, e := range data {
			err = recv(ctx, e)
			if err != nil {
				return
			}
		}

		return
	})
}

func Map[T, E any](it Iterable[T], mapper func(ctx context.Context, data T) (E, error)) Iterable[E] {
	return IteratorFunc[E](func(ctx context.Context, res Receiver[E]) (err error) {
		return it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
			mapped, err := mapper(ctx, data)
			if err != nil {
				return
			}
			return res(ctx, mapped)
		}))
	})
}
