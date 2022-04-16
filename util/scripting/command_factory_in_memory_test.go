package scripting_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/teawithsand/handmd/util/scripting"
)

func TestInMemoryCommandFactory(t *testing.T) {
	ctx := context.Background()

	fac := scripting.InMemoryCommandFactory{
		Scripts: scripting.InMemoryScripts{
			"test.sh": scripting.InMemoryScript{
				Script: `#!/bin/bash
echo "it works!";
`,
			},
		},
	}
	defer fac.Close()

	err := fac.Initialize()
	if err != nil {
		t.Error(err)
		return
	}

	cmd, err := fac.GetCommand(ctx, "test.sh")
	if err != nil {
		t.Error(err)
		return
	}

	b := bytes.NewBuffer(nil)
	cmd.Stdout = b

	err = cmd.Exec(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	if b.String() != "it works!\n" {
		t.Error("got invalid data", b.String())
		return
	}
}
