package renderer

import (
	"context"
	"encoding/json"
	"os"

	"github.com/teawithsand/handmd/util/encoding"
	"github.com/teawithsand/handmd/util/fsal"
)

// Encoding renders single entry of type T to specified file with specified format.
type Encoding[T any] struct {
	EncoderFactory encoding.EncoderFactory
	FileName       string
}

func (er *Encoding[T]) Render(ctx context.Context, input T, fs RendererOutput) (err error) {
	var f fsal.File
	f, err = fs.Open(er.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return
	}
	defer f.Close()

	var enc encoding.Encoder
	if er.EncoderFactory != nil {
		enc = er.EncoderFactory.MakeEncoder(f)
	} else {
		innerEnc := json.NewEncoder(f)
		innerEnc.SetIndent("", "\t")
		enc = innerEnc
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
