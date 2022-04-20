// Package ctn contains utils for making DI containers.
// It wasn't crated as general-purpose DI, but something to make compilation process simpler.
package ctn

import (
	"go.uber.org/dig"
)

type Container struct {
	dig *dig.Container
}

// HACK(teawithsand): for now do not implement proper container, just ust
func (c *Container) Inner() *dig.Container {
	return c.dig
}

func (c *Container) Close() (err error) {
	err = c.dig.Invoke(func(cln Cleaner) (err error) {
		err = cln.Close()
		return
	})
	return
}

// Container, created especially for compilation process.
func NewCompilationContainer() (ctn *Container, err error) {
	c := dig.New()
	err = c.Provide(func() (reg Cleaner, err error) {
		reg = NewCleaner()
		return
	})
	if err != nil {
		return
	}

	err = c.Provide(func(cln Cleaner) (reg CleanRegistry, err error) {
		reg = cln
		return
	})
	if err != nil {
		return
	}

	ctn = &Container{
		dig: c,
	}
	return
}
