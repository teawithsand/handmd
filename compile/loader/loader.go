package loader

import (
	"context"

	"github.com/teawithsand/handmd/util/fsal"
)

type LoaderInput = fsal.FS

type Loader[C any, O any] interface {
	Load(ctx context.Context, lctx C, src LoaderInput) (res O, err error)
}

type LoaderFunc[C any, O any] func(ctx context.Context, lctx C, input LoaderInput) (output O, err error)

func (f LoaderFunc[C, O]) Load(ctx context.Context, lctx C, input LoaderInput) (output O, err error) {
	return f(ctx, lctx, input)
}

var _ Loader[any, any] = LoaderFunc[any, any](func(ctx context.Context, lctx any, input LoaderInput) (output any, err error) {
	return
})
