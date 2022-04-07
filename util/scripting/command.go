package scripting

import (
	"context"
	"io"
	"os/exec"
)

// Wrapper for go exec.Command type, which is simpler.
// Adds support for context.
type Command struct {
	Command string
	Args    []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Command) Exec(ctx context.Context) (err error) {
	cmd := exec.Command(c.Command, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stdin = c.Stdin
	cmd.Stderr = c.Stderr

	doneChan := make(chan error)
	go func() {
		doneChan <- cmd.Run()
	}()

	select {
	case <-ctx.Done():
		cmd.Process.Kill()

		err = ctx.Err()
		return
	case err = <-doneChan:
		if err != nil {
			return
		}
	}

	return
}
