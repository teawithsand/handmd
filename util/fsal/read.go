package fsal

import (
	"bytes"
	"errors"
	"io"
	"os"
)

// Note: this function is inherently unsafe, since it does not limit incoming file size.
// Use only in trusted environments and with caution.
func ReadFile(fs FS, path string) (data []byte, err error) {
	b := bytes.NewBuffer(nil)

	f, err := fs.Open(path, os.O_RDONLY)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(b, f)
	if err != nil {
		return
	}

	data = b.Bytes()

	return
}

type DataParser struct {
	Parser    func(data []byte, res interface{}) (err error)
	Extension string

	// MaxInputSize int64 // NIY
}

// Util for loading data from file in multiple formats.
// Internally uses ReadFile, so it's not safe against big files or reading /dev/urandom.
type DataLoader struct {
	Parsers []DataParser
}

func (dl *DataLoader) ReadData(fs FS, path string, res interface{}) (err error) {
	for _, p := range dl.Parsers {
		var data []byte
		data, err = ReadFile(fs, path+"."+p.Extension)
		if errors.Is(err, os.ErrNotExist) {
			err = nil
			continue
		}

		err = p.Parser(data, res)
		if err != nil {
			return
		}
		return
	}

	err = os.ErrNotExist
	return
}
