package fsal

import (
	"io/ioutil"
	"os"
)

type LocalFS struct {
}

var _ FS = &LocalFS{}

func (fs *LocalFS) Open(path string, openMode int) (f File, err error) {
	f, err = os.OpenFile(path, openMode, 0660)
	return
}
func (fs *LocalFS) ReadDir(path string) (entries []Entry, err error) {
	rawEntries, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, re := range rawEntries {
		re := re
		entries = append(entries, re)
	}
	return
}

func (fs *LocalFS) Rename(from, to string) (err error) {
	return os.Rename(from, to)
}

func (fs *LocalFS) Mkdir(path string) (err error) {
	err = os.Mkdir(path, 0770)
	return
}

func (fs *LocalFS) MkdirAll(path string) (err error) {
	err = os.MkdirAll(path, 0770)
	return
}

func (fs *LocalFS) Remove(path string) (err error) {
	err = os.Remove(path)
	return
}

func (fs *LocalFS) RemoveAll(path string) (err error) {
	if path == "/" {
		panic("no no, this is not ok")
	}
	err = os.RemoveAll(path)
	return
}

func (fs *LocalFS) Stat(path string) (entry Entry, err error) {
	return os.Stat(path)
}
