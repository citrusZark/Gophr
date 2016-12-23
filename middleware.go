package main

import "net/http"

type Middleware []http.Handler

type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(bytes)
}

func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func NewMiddlewareResponseWriter(w http.ResponseWriter) *MiddlewareResponseWriter {
	return &MiddlewareResponseWriter{
		ResponseWriter: w,
	}
}

// Adds a handler to Middleware
func (m *Middleware) Add(handler http.Handler) {
	*m = append(*m, handler)
}
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Wrap the supplied ResponseWriter
	mw := NewMiddlewareResponseWriter(w)

	// Loop through all the registered handler
	for _, handler := range m {
		// Call the handler with our MiddlewareResponseWriter
		handler.ServeHTTP(mw, r)

		// If there was a write, stop processing
		if mw.written {
			return
		}
	}

	// If no handler wrote response, it's 404
	http.NotFound(w, r)
}
