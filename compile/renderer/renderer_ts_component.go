package renderer

import (
	"context"
	"os"

	"github.com/teawithsand/handmd/util/tsrender"
)

// Simple renderer for rendering typescript components.
// Uses bunch of hacks that I've created.
// Feel free to write your own renderer.
type TSComponent struct {
	TargetPath  string
	BaseImports []tsrender.Import
}

type TSComponentData struct {
	Impots []tsrender.Import
}

type TSRednerData struct {
	Impots []tsrender.Import
	Tag    tsrender.SimpleTag
}

func (tcr *TSComponent) Render(ctx context.Context, tag tsrender.SimpleTag, fs RendererOutput) (err error) {
	f, err := fs.Open(tcr.TargetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return
	}
	defer f.Close()

	for _, imp := range tcr.BaseImports {
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
		`export default (props: any) => {
return (`,
	))
	if err != nil {
		return
	}

	err = tag.Render(ctx, f)
	if err != nil {
		return
	}

	_, err = f.Write([]byte(
		`)
}`,
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
