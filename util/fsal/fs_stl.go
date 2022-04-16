package fsal

import (
	"io/fs"
	"os"
)

// Wraps fsal file system and provides "io/fs" FS from it.
// It also makes FSAL compatible with http.FileSystem via http.FS function.
type STLFS struct {
	Wrapped FS
}

var _ fs.FS = &STLFS{}

func (fs *STLFS) Open(name string) (f fs.File, err error) {
	f, err = fs.Wrapped.Open(name, os.O_RDONLY)
	return
}
