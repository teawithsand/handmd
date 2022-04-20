package fsal_test

import (
	"testing"

	"github.com/teawithsand/handmd/util/fsal"
)

func TestTempDir(t *testing.T) {
	td, err := fsal.NewTempDir(nil, "asdf_")
	if err != nil {
		t.Error(err)
		return
	}
	defer td.Close()
}
