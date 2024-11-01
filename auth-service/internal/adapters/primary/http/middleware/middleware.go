package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type chain struct {
	middlewares []Middleware
}

func NewChain(middlewares ...Middleware) *chain {
	return &chain{
		middlewares: append([]Middleware{}, middlewares...),
	}
}

func (c *chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}

	return h
}

func (c *chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}
