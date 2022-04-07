package httpext

import "net/http"

type Middleware interface {
	Apply(h http.Handler) http.Handler
}

type MiddlewareFunc func(h http.Handler) http.Handler

func (f MiddlewareFunc) Apply(h http.Handler) http.Handler {
	return f(h)
}

type Middlewares []Middleware

func (mws Middlewares) Apply(h http.Handler) http.Handler {
	for _, mw := range mws {
		h = mw.Apply(h)
	}

	return h
}
