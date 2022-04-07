package renderer

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/teawithsand/handmd/util/fsal"
	"github.com/teawithsand/handmd/util/scripting"
)

// Renderer, which creates fuse.js indexes using fuseIndex.js from commonscripts.
type FuseIndex[T any] struct {
	CommandFactory scripting.CommandFactory
	FuseComandName string

	FuseIndexOutputPath string
	IndexFields         []string
}

func (dr *FuseIndex[T]) Render(ctx context.Context, input []T, fs RendererOutput) (err error) {
	var f fsal.File
	f, err = fs.Open(dr.FuseIndexOutputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return
	}
	defer f.Close()

	var cmd *scripting.Command
	cmd, err = dr.CommandFactory.GetCommand(ctx, dr.FuseComandName)
	if err != nil {
		return
	}

	cmd.Args = dr.IndexFields

	pr, pw := io.Pipe()
	wg := sync.WaitGroup{}
	wg.Add(1)

	var encodeError error
	go func() {
		defer wg.Done()
		defer pw.Close()

		enc := json.NewEncoder(pw)
		err = enc.Encode(input)
		if encodeError != nil {
			return
		}
	}()

	cmd.Stdin = pr
	cmd.Stdout = f

	err = cmd.Exec(ctx)
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	_, err = io.Copy(io.Discard, pr)
	if err != nil {
		return
	}

	wg.Wait()

	if encodeError != nil {
		err = encodeError
		return
	}

	return
}
