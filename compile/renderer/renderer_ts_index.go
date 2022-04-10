package renderer

import (
	"context"
	"fmt"
	"os"

	"github.com/teawithsand/handmd/util/iter"
	"github.com/teawithsand/handmd/util/tsrender"
)

type TSIndex struct {
	BaseImports []tsrender.Import
	TargetPath  string
}

type TSIndexRenderEntry = map[string]string

type TSIndexRenderData struct {
	Imports []tsrender.Import

	// Map of string value, to string, value, which will be pasted as raw string into typescript.
	Entries iter.Iterable[TSIndexRenderEntry]
}

func (tir *TSIndex) Render(ctx context.Context, data TSIndexRenderData, fs RendererOutput) (err error) {
	f, err := fs.Open(tir.TargetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return
	}
	defer f.Close()

	for _, imp := range tir.BaseImports {
		_, err = f.Write([]byte(imp.Render()))
		if err != nil {
			return
		}

		_, err = f.Write([]byte("\n"))
		if err != nil {
			return
		}
	}

	for _, imp := range data.Imports {
		_, err = f.Write([]byte(imp.Render()))
		if err != nil {
			return
		}

		_, err = f.Write([]byte("\n"))
		if err != nil {
			return
		}
	}

	_, err = f.Write([]byte(
		"export default [\n",
	))
	if err != nil {
		return
	}

	err = data.Entries.Iterate(ctx, iter.Receiver[TSIndexRenderEntry](func(ctx context.Context, data TSIndexRenderEntry) (err error) {
		_, err = f.Write([]byte("{\n"))
		if err != nil {
			return
		}

		for k, v := range data {
			_, err = f.Write([]byte(fmt.Sprintf("\"%s\": %s,", k, v)))
			if err != nil {
				return
			}
		}

		_, err = f.Write([]byte("},\n"))
		if err != nil {
			return
		}
		return
	}))
	if err != nil {
		return
	}

	_, err = f.Write([]byte(
		`]`,
	))
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	return
}
