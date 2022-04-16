package scripting

import (
	"context"
	"errors"
	"os"
	"path"
)

type InMemoryScript struct {
	Script string
}

type InMemoryScripts map[string]InMemoryScript

// Note: init and close methods of this factory are not goroutine-safe.
type InMemoryCommandFactory struct {
	Scripts InMemoryScripts
	dir     string
}

// Stores scripts in temporary directory.
// Note: this method along with Close is not thread-safe.
func (fac *InMemoryCommandFactory) Initialize() (err error) {
	dir, err := os.MkdirTemp(os.TempDir(), "handmd_scripts_*")
	if err != nil {
		return
	}

	for scriptName, script := range fac.Scripts {
		// TODO(teawithsand): sanitize scriptName
		p := path.Join(dir, scriptName)

		err = os.WriteFile(p, []byte(script.Script), 0550)
		if err != nil {
			return
		}
	}

	fac.dir = dir
	return
}

func (fac *InMemoryCommandFactory) GetCommand(ctx context.Context, name string) (cmd *Command, err error) {
	if len(fac.dir) == 0 {
		err = errors.New("handmd/util/scripting: InMemoryScriptCommandFactory was not initialized")
		return
	}

	cmd = &Command{
		Command: path.Join(fac.dir, name),
	}

	return
}

// Note: this method along with Initialize is not thread-safe.
func (fac *InMemoryCommandFactory) Close() (err error) {
	if len(fac.dir) == 0 {
		return nil
	}

	d := fac.dir

	fac.dir = ""
	err = os.RemoveAll(d)
	if err != nil {
		return
	}

	return
}
