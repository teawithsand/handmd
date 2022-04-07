package renderer

import (
	"context"
	"fmt"

	"github.com/teawithsand/handmd/util/fsal"
)

type Copy struct {
	SourcePath      string
	DestinationPath string
	Required        bool
}

func (cr *Copy) Render(ctx context.Context, srcFs fsal.FS, dstFs RendererOutput) (err error) {
	exists, err := fsal.Exists(srcFs, cr.SourcePath)
	if err != nil {
		return
	}

	if !exists && cr.Required {
		err = fmt.Errorf("can't copy %s; it does not exist", cr.SourcePath)
		return
	}
	if !exists {
		return
	}

	err = dstFs.MkdirAll(cr.DestinationPath)
	if err != nil {
		return
	}

	err = fsal.CopyDirectory(srcFs, dstFs, cr.SourcePath, cr.DestinationPath)
	if err != nil {
		return
	}
	return
}
