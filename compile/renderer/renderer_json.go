package renderer

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/teawithsand/handmd/util/fsal"
)

// JSON renders single entry of type T to specfied json file.
type JSON[T any] struct {
	EncoderFactory func(w io.Writer) *json.Encoder
	FileName       string
}

func (jr *JSON[T]) Render(ctx context.Context, input T, fs RendererOutput) (err error) {
	var f fsal.File
	f, err = fs.Open(jr.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return
	}
	defer f.Close()

	var enc *json.Encoder
	if jr.EncoderFactory != nil {
		enc = jr.EncoderFactory(f)
	} else {
		enc = json.NewEncoder(f)
		enc.SetIndent("", "\t")
	}
	err = enc.Encode(input)
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}
	return
}
