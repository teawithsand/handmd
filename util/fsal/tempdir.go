package fsal

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type TempDir struct {
	fs   FS
	name string
}

func (td *TempDir) Close() (err error) {
	// use once for this purpose?
	err = td.FS().RemoveAll(".")
	if err != nil {
		return
	}
	return
}

// Returns FS prefixed with directory created
func (td *TempDir) FS() FS {
	return &PrefixFS{
		Wrapped:    td.fs,
		PathPrefix: td.name,
	}
}

// Returns directory name.
func (td *TempDir) Name() string {
	return td.name
}

// Returns FS that this temp dir was created in.
func (td *TempDir) ParentFS() FS {
	return td.fs
}

// Creates temp dir in specified fs.
// Removes it once temp dir is closed.
//
// If provided FS is nil, LocalTempFS() is used.
func NewTempDir(fs FS, prefix string) (dir *TempDir, err error) {
	if fs == nil {
		fs = LocalTempFS()
	}

	var buf [8]byte
	_, err = io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return
	}
	name := fmt.Sprintf("%s%s", prefix, hex.EncodeToString(buf[:]))

	err = fs.Mkdir(name)
	if err != nil {
		return
	}

	dir = &TempDir{
		fs:   fs,
		name: name,
	}
	return
}
