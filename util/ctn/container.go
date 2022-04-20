// Package ctn contains utils for making DI containers.
// It wasn't crated as general-purpose DI, but something to make compilation process simpler.
package ctn

type Container struct {
	dig DI
}

// HACK(teawithsand): for now do not implement proper container, just ust
func (c *Container) Inner() DI {
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
	c := NewDI()
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