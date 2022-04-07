package fsal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

// TODO(teawithsand): inter filesystem copy

func CopyDirectory(srcFs, dstFs FS, scrDir, dest string) error {
	if dstFs == nil {
		dstFs = srcFs
	}

	entries, err := srcFs.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := path.Join(scrDir, entry.Name())
		destPath := path.Join(dest, entry.Name())

		fileInfo, err := srcFs.Stat(sourcePath)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if err := createIfNotExists(dstFs, destPath); err != nil {
				return err
			}
			if err := CopyDirectory(srcFs, dstFs, sourcePath, destPath); err != nil {
				return err
			}
		} else {
			if err := innerCopy(srcFs, dstFs, sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func innerCopy(srcFs, dstFs FS, srcFile, dstFile string) error {
	if dstFs == nil {
		dstFs = srcFs
	}

	out, err := dstFs.Open(dstFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := srcFs.Open(srcFile, os.O_RDONLY)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	return nil
}

func Exists(fs FS, filePath string) (ok bool, err error) {
	if _, err := fs.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return err == nil, err
}

func createIfNotExists(fs FS, dir string) error {
	exists, err := Exists(fs, dir)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if err := fs.MkdirAll(dir); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
