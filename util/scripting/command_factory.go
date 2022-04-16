package scripting

import (
	"context"
	"path"
)

// CommandFactory creates commands from their name.
type CommandFactory interface {
	// Returned command has only path filled. Rest is up to user.
	GetCommand(ctx context.Context, name string) (cmd *Command, err error)
}

type FSCommandFactory struct {
	Dir string
}

func (fac *FSCommandFactory) GetCommand(ctx context.Context, name string) (cmd *Command, err error) {
	cmd = &Command{
		Command: path.Join(fac.Dir, name),
	}
	return
}
