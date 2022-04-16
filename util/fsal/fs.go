package fsal

import (
	"io"
	"io/fs"
)

type FileInfo = fs.FileInfo

type File interface {
	io.Closer
	io.Reader
	io.Writer
	io.Seeker
	io.ReaderAt
	io.WriterAt

	Stat() (FileInfo, error)
}

type Entry interface {
	IsDir() bool
	Name() string
}

type FS interface {
	Open(path string, openMode int) (f File, err error)

	ReadDir(path string) (entries []Entry, err error)
	Rename(from, to string) (err error)
	Mkdir(path string) (err error)
	Remove(path string) (err error) // also works as RMDIr
	Stat(path string) (e FileInfo, err error)

	RemoveAll(path string) (err error)
	MkdirAll(path string) (err error)
}
