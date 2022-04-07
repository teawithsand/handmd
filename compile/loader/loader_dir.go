package loader

import (
	"context"
	"strings"

	"github.com/teawithsand/handmd/util/fsal"
	"github.com/teawithsand/handmd/util/iter"
)

type DirWrapper[C any] struct {
	LoaderContext C
	DirName       string
}

func (pw *DirWrapper[C]) Unwrap() C {
	return pw.LoaderContext
}

// Loader whcih loads defines.RawPost
type Dir[C any, T any] struct {
	InnerLoader Loader[DirWrapper[C], T]

	Filter func(ctx context.Context, name string) (ignore bool, err error) // defaults to ignoring dirs starting with . and _
}

func (dl *Dir[C, T]) Load(ctx context.Context, lctx C, fs LoaderInput) (iterator iter.Iterable[T], err error) {
	iterator = iter.IterableFunc[T](func(ctx context.Context, recv iter.Receiver[T]) (err error) {
		entries, err := fs.ReadDir("/")
		if err != nil {
			return
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			postDir := e.Name()
			if dl.Filter == nil {
				if strings.HasPrefix(postDir, ".") || strings.HasPrefix(postDir, "_") {
					continue
				}
			} else {
				var ignore bool
				ignore, err = dl.Filter(ctx, postDir)
				if err != nil {
					return
				}
				if ignore {
					continue
				}
			}

			ifs := &fsal.PrefixFS{
				Wrapped:    fs,
				PathPrefix: postDir,
			}

			var res T
			res, err = dl.InnerLoader.Load(ctx, DirWrapper[C]{
				LoaderContext: lctx,
				DirName:       postDir,
			}, ifs)
			if err != nil {
				return
			}

			err = recv(ctx, res)
			if err != nil {
				return
			}
		}

		return
	})

	return
}
