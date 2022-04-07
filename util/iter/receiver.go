package iter

import "context"

// Function, which is used to get data from iterable.
type Receiver[T any] func(ctx context.Context, data T) (err error)
