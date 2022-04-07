package renderer

import (
	"context"
)

// Renderer, which removes EVERTTHING IN PROVIDED OUTPUT!!!!
// It does rm -rf / ON GIVEN FS.
// NEVER GIVE IT LOCAL FS.
type Clean struct {
}

// It does rm -rf / ON GIVEN OUTPUT FS.
// !!!!NEVER GIVE IT LOCAL FS!!!1!!
func (jr *Clean) Render(ctx context.Context, input struct{}, fs RendererOutput) (err error) {
	err = fs.RemoveAll("/")
	if err != nil {
		return
	}

	err = fs.Mkdir("/")
	if err != nil {
		return
	}
	return
}
