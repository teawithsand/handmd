package iter

import (
	"context"
)

// Collection of elements of type T, which may be iterated on.
// Note: iterables are allowed to be concurrent. They may launch multiple goroutines at the same time, even though that's not the default.
type Iterable[T any] interface {
	Iterate(ctx context.Context, recv Receiver[T]) (err error)
}

type IterableFunc[T any] func(ctx context.Context, recv Receiver[T]) (err error)

func (f IterableFunc[T]) Iterate(ctx context.Context, recv Receiver[T]) (err error) {
	return f(ctx, recv)
}

// Sends arr elements from iterable given to specified target channel.
// Blocks current goroutine.
func IterateToChannel[T any](ctx context.Context, it Iterable[T], target chan<- T) (err error) {
	err = it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
		select {
		case target <- data:
			return
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}))
	return
}

func Collect[T any](ctx context.Context, it Iterable[T]) (res []T, err error) {
	err = it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
		res = append(res, data)
		return
	}))
	return
}

func Slice[T any](data []T) Iterable[T] {
	return IterableFunc[T](func(ctx context.Context, recv Receiver[T]) (err error) {
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
	return IterableFunc[E](func(ctx context.Context, res Receiver[E]) (err error) {
		return it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
			mapped, err := mapper(ctx, data)
			if err != nil {
				return
			}
			return res(ctx, mapped)
		}))
	})
}

func Filter[T any](it Iterable[T], filter func(ctx context.Context, data T) (bool, error)) Iterable[T] {
	return IterableFunc[T](func(ctx context.Context, res Receiver[T]) (err error) {
		return it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
			ok, err := filter(ctx, data)
			if err != nil {
				return
			}
			if ok {
				return res(ctx, data)
			}
			return
		}))
	})
}

// Note: iterator given must not be parallel one.
func JoinString(ctx context.Context, it Iterable[string], sep string) (res string, err error) {
	isFirst := true
	err = it.Iterate(ctx, Receiver[string](func(ctx context.Context, data string) (err error) {
		if !isFirst {
			res += "; "
		}
		isFirst = false

		res += data
		return
	}))

	return
}

/*
This function is too complex and I am really not sure if it should be part of handmd.
For sake of simplicity, stored as comment rather than branch or sth.

type closerFunc func() (err error)

func (f closerFunc) Close() (err error) {
	return f()
}

var ErrIteratorClosed = errors.New("handmd/util/iter: iterator was closed before end of iteration")

// Creates channel, which receives elements from iterable given.
// Iterable will be called in separate goroutine.
func IterateAsChannel[T any](ctx context.Context, it Iterable[T], done chan<- error) (target chan<- T, closer func()) {
	closeChan := make(chan struct{})
	once := sync.Once{}
	closer = func() {
		once.Do(func() {
			close(closeChan)
		})
	}

	go func() {
		done <- it.Iterate(ctx, Receiver[T](func(ctx context.Context, data T) (err error) {
			select {
			case target <- data:
				return
			case <-closeChan:
				err = ErrIteratorClosed
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
		}))
	}()
	return
}*/
