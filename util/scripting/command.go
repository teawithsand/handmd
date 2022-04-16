package scripting

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync/atomic"
)

// Wrapper for go exec.Command type, which is simpler.
// Adds support for context.
type Command struct {
	Command string
	Args    []string
	EnvVars map[string]string

	// If empty, current process' CWD is used
	Dir string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Command) Exec(ctx context.Context) (err error) {
	cmd := exec.CommandContext(ctx, c.Command, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stdin = c.Stdin
	cmd.Stderr = c.Stderr
	cmd.Dir = c.Dir

	for k, v := range c.EnvVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	err = cmd.Run()
	return
}

// Executes command until:
// 1. Context gets cancelled
// 2. close function gets called.
//
// Once command gets closed, result error is returned through result channel.
// If command was closed via returned closer, nil error is emitted.
// Panics, if channel is closed.
//
// Not closing command causes goroutine leak, so it should be always closed either
// via canceling context or calling closer function.
func (c *Command) ExecBackground(ctx context.Context, res chan<- error) (close func()) {
	var flag int32 // 1 once result was sent via channel

	ctx, canceler := context.WithCancel(ctx)
	go func() {
		err := c.Exec(ctx)
		if atomic.CompareAndSwapInt32(&flag, 0, 1) {
			res <- err
		}
	}()

	close = func() {
		if atomic.CompareAndSwapInt32(&flag, 0, 1) {
			res <- nil
			canceler()
		}
	}

	return
}
