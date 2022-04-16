package fsal

import (
	"embed"
	"io"
	"io/fs"
)

type EmbedFS struct {
	FS embed.FS
}

var _ FS = &EmbedFS{}

type passedFile interface {
	fs.File
	io.Seeker
}

type fileWrapper struct {
	passedFile
}

func (fw *fileWrapper) Close() (err error) {
	return
}

func (fw *fileWrapper) Write(data []byte) (n int, err error) {
	err = ErrReadOnlyFS
	return
}

func (fw *fileWrapper) WriteAt(data []byte, off int64) (n int, err error) {
	err = ErrReadOnlyFS
	return
}

func (fw *fileWrapper) ReadAt(data []byte, off int64) (n int, err error) {
	// HACK: this should be implemented properly
	// it's bypass, since readAt is hardly ever used
	err = ErrOperationNotSupported
	return
}

func (fs *EmbedFS) Open(path string, openMode int) (f File, err error) {
	rawFile, err := fs.FS.Open(path)
	if err != nil {
		return
	}
	f = &fileWrapper{
		passedFile: rawFile.(passedFile),
	}
	return
}

func (fs *EmbedFS) ReadDir(path string) (entries []Entry, err error) {
	rawEntries, err := fs.FS.ReadDir(path)
	if err != nil {
		return
	}
	for _, re := range rawEntries {
		re := re
		entries = append(entries, re)
	}
	return
}

func (fs *EmbedFS) Rename(from, to string) (err error) {
	err = ErrReadOnlyFS
	return
}

func (fs *EmbedFS) Mkdir(path string) (err error) {
	err = ErrReadOnlyFS
	return
}

func (fs *EmbedFS) MkdirAll(path string) (err error) {
	err = ErrReadOnlyFS
	return
}

func (fs *EmbedFS) Remove(path string) (err error) {
	err = ErrReadOnlyFS
	return
}

func (fs *EmbedFS) RemoveAll(path string) (err error) {
	err = ErrReadOnlyFS
	return
}

func (fs *EmbedFS) Stat(inputPath string) (entry FileInfo, err error) {
	f, err := fs.FS.Open(inputPath)
	if err != nil {
		return
	}
	entry, err = f.Stat()
	if err != nil {
		return
	}
	return
}
