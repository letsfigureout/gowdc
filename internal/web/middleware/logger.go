package middleware

import (
	"log"
	"net/http"
	"time"
)

// Middleware struct
type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() {
				log.Printf("%s %s -> %s (%s)",
					r.Method,
					r.URL.RequestURI(),
					r.RemoteAddr,
					time.Since(start))
			}()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, mw := range middlewares {
		f = mw(f)
	}
	return f
}
