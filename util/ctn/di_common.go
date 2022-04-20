package ctn

import (
	"io"

	"github.com/teawithsand/handmd/util/fsal"
)

// Common DI type for input file system in DI container.
type InputFS fsal.FS

// Common DI type for output file system in DI container.
type OutputFS fsal.FS

// Initializes DI, so that DICommonUtil can be used with it.
// For now it only registers `Cleaner` and `CleanRegistry`.
func InitializeDIForUtil(di DI) (err error) {
	err = di.Provide(func() (cleaner Cleaner, err error) {
		cleaner = NewCleaner()
		return
	})
	if err != nil {
		return
	}

	err = di.Provide(func(cleaner Cleaner) (reg CleanRegistry, err error) {
		reg = cleaner
		return
	})
	if err != nil {
		return
	}

	return
}

// Getters/accessors for common utils used in DI by compilation stuff.
// Note: di must be initialized first in order to use it.
type DICommonUtil struct {
	DI

	// TODO(teawithsand): add some cache here, since invoke calls may be somewhat slow
	// and all values are singletons anyway, so using cache is ok
}

func (di *DICommonUtil) RegisterCloser(c io.Closer) (err error) {
	return di.DI.Invoke(func(reg CleanRegistry) (err error) {
		reg.Register(c)
		return
	})
}

// Cleans up all closes, which were registered in CleanRegistry.
func (di *DICommonUtil) CleanupClosers() (err error) {
	return di.DI.Invoke(func(reg Cleaner) (err error) {
		err = reg.Clean()
		return
	})
}
