package fsal

import (
	"path"
	"strings"
)

// PrefixFS prefixes all paths with prefix and does not allow going above it.
// Works like chroot without bypass.
type PrefixFS struct {
	Wrapped    FS
	PathPrefix string
}

var _ FS = &PrefixFS{}

func (fs *PrefixFS) translatePath(inputPath string) string {
	pp := path.Clean(fs.PathPrefix)

	newPath := path.Clean(path.Join(pp, inputPath))
	if !strings.HasPrefix(newPath, pp) {
		return fs.PathPrefix
	}
	if strings.Contains(newPath, "..") {
		panic("result path is relative! It should not happen")
	}

	return newPath
}

func (fs *PrefixFS) Open(path string, openMode int) (f File, err error) {
	return fs.Wrapped.Open(fs.translatePath(path), openMode)
}

func (fs *PrefixFS) ReadDir(path string) (entries []Entry, err error) {
	return fs.Wrapped.ReadDir(fs.translatePath(path))
}

func (fs *PrefixFS) Rename(from, to string) (err error) {
	return fs.Wrapped.Rename(fs.translatePath(from), fs.translatePath(to))
}

func (fs *PrefixFS) Mkdir(path string) (err error) {
	err = fs.Wrapped.Mkdir(fs.translatePath(path))
	return
}

func (fs *PrefixFS) MkdirAll(path string) (err error) {
	err = fs.Wrapped.MkdirAll(fs.translatePath(path))
	return
}

func (fs *PrefixFS) Remove(path string) (err error) {
	err = fs.Wrapped.Remove(fs.translatePath(path))
	return
}

func (fs *PrefixFS) RemoveAll(path string) (err error) {
	err = fs.Wrapped.RemoveAll(fs.translatePath(path))
	return
}

func (fs *PrefixFS) Stat(path string) (entry FileInfo, err error) {
	return fs.Wrapped.Stat(fs.translatePath(path))
}
