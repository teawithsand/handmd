package ctn

import "go.uber.org/dig"

// Implementation of DI container used.
// For now dig is used, rather than some interface for sake of simplicity.
type DI = *dig.Container

// Creates new empty DI container.
func NewDI() DI {
	return dig.New()
}
