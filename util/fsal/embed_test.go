package fsal_test

import (
	"bytes"
	"embed"
	"io"
	"os"
	"testing"

	"github.com/teawithsand/handmd/util/fsal"
)

//go:embed *.go
var data embed.FS

func TestEmbed_CanReadFile(t *testing.T) {
	efs := fsal.EmbedFS{
		FS: data,
	}

	f, err := efs.Open("embed_test.go", os.O_RDONLY)
	if err != nil {
		t.Error(err)
		return
	}

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, f)
	if err != nil {
		t.Error(err)
		return
	}

	if b.Len() == 0 {
		t.Error("got no data")
		return
	}
}

func TestEmbed_CanStat(t *testing.T) {
	efs := fsal.EmbedFS{
		FS: data,
	}

	f, err := efs.Stat("embed_test.go")
	if err != nil {
		t.Error(err)
		return
	}

	_ = f
}
