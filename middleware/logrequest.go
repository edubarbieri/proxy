package middleware

import (
	"log"
	"net/http"
	"time"
)

type LogMiddleware struct {
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

}
func (md *LogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(res)
		next.ServeHTTP(wrapped, req)
		log.Printf("| %3d | %13v | %v |%s %-7s", wrapped.status, time.Since(start), GetRemoteIP(req), req.Method, req.URL.EscapedPath())
	})
}
