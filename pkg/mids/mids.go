package mids

import (
	"net/http"
)

// Middleware is a middleware.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// New creates a new Middleware.
// mid1(mid2(mid3(handlerFunc))) is equivalent to
// New(handlerFunc)(mid1, mid2, mid3)
func New(next http.HandlerFunc) func(ms ...Middleware) http.HandlerFunc {
	return func(ms ...Middleware) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			res := next

			for i := len(ms) - 1; i >= 0; i-- {
				res = ms[i](res)
			}

			res(w, r)
		}
	}
}
