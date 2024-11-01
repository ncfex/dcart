package middleware

import (
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		method := r.Method
		fmt.Printf("Testing logging middleware:\nMethod:%s\nUser-Agent:%s\n\n", method, userAgent)
		next.ServeHTTP(w, r)
	}
}
