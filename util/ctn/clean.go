package ctn

import (
	"io"
	"sync"
)

type CleanRegistry interface {
	Register(closeable io.Closer)
}

type Cleaner interface {
	io.Closer
	CleanRegistry
}

type cleaner struct {
	lock    sync.Mutex
	closers []io.Closer
}

func (c *cleaner) Register(closer io.Closer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.closers = append(c.closers, closer)
}

func (c *cleaner) Close() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, e := range c.closers {
		tempErr := e.Close()
		if tempErr != nil {
			err = tempErr
		}
	}

	return
}

func NewCleaner() Cleaner {
	return &cleaner{}
}
