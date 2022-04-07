package renderer

import (
	"context"

	"github.com/teawithsand/handmd/util/fsal"
	"github.com/teawithsand/handmd/util/iter"
)

type DirWrapper[T any] struct {
	Data T

	Dir          string
	ParentOutput RendererOutput
}

func (pw *DirWrapper[T]) Unwrap() T {
	return pw.Data
}

type Dir[T any] struct {
	// Function, which for given data picks directory for each incoming data.
	DirPicker func(ctx context.Context, data T) (dir string, err error)

	// Renderer, which will be invoked for each post output dir.
	PostRenderer Renderer[DirWrapper[T]]
}

func (pr *Dir[T]) Render(ctx context.Context, input iter.Iterable[T], fs RendererOutput) (err error) {
	return input.Iterate(ctx, iter.Receiver[T](func(ctx context.Context, data T) (err error) {
		dir, err := pr.DirPicker(ctx, data)
		if err != nil {
			return
		}

		postFs := &fsal.PrefixFS{
			Wrapped:    fs,
			PathPrefix: dir,
		}

		err = postFs.Mkdir("/")
		if err != nil {
			return
		}

		err = pr.PostRenderer.Render(ctx, DirWrapper[T]{
			Dir:          dir,
			Data:         data,
			ParentOutput: fs,
		}, postFs)
		if err != nil {
			return
		}

		return
	}))
}
