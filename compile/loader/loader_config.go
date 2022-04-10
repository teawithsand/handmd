package loader

import (
	"context"

	"github.com/teawithsand/handmd/util/fsal"
)

type Config[C any, O any] struct {
	Loader   fsal.DataLoader
	Factory  func() *O // must always return non-nil pointer
	FileName string
}

var _ Loader[any, any] = &Config[any, any]{}

func (cl *Config[C, O]) Load(ctx context.Context, lctx C, fs LoaderInput) (config O, err error) {
	v := cl.Factory()
	err = cl.Loader.ReadData(fs, cl.FileName, v)
	if err != nil {
		return
	}
	config = *v

	return
}
