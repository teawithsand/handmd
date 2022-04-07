package fsal

import (
	"bytes"
	"io"
	"os"
)

func WriteFile(fs FS, path string, data []byte) (err error) {
	return WriteFileStream(fs, path, bytes.NewReader(data))
}

func WriteFileStream(fs FS, path string, r io.Reader) (err error) {
	f, err := fs.Open(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return
	}

	// sync? For most use cases, it's not needed.

	err = f.Close()
	if err != nil {
		return
	}

	return
}
