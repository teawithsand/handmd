package renderer

import (
	"context"

	"github.com/teawithsand/handmd/util/fsal"
)

// Destination of renderer rendering.
// For now it's always FS.
// In future this type may change.
type RendererOutput = fsal.FS

// Render takes some data(and configuration) and performs rendering into RendererOutput.
//
// Typically, renderers are split into two categories: global ones and
// post ones, depending on output location they use.
type Renderer[I any] interface {
	Render(ctx context.Context, input I, output RendererOutput) (err error)
}

type RendererFunc[I any] func(ctx context.Context, input I, output RendererOutput) (err error)

func (f RendererFunc[I]) Render(ctx context.Context, input I, output RendererOutput) (err error) {
	return f(ctx, input, output)
}
