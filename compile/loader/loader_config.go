package loader

import (
	"context"

	"github.com/teawithsand/handmd/util/fsal"
)

type Config[T any] struct {
	Loader   fsal.DataLoader
	Factory  func() *T // must always return non-nil pointer
	FileName string
}

func (cl *Config[T]) Load(ctx context.Context, fs LoaderInput) (config T, err error) {
	v := cl.Factory()
	err = cl.Loader.ReadData(fs, cl.FileName, v)
	if err != nil {
		return
	}
	config = *v

	return
}
