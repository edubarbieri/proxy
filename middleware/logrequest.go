package middleware

import (
	"log"
	"net/http"
	"time"
)

type LogMiddleware struct {
}

func (md *LogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(res)
		next.ServeHTTP(wrapped, req)
		log.Printf("| %3d | %13v | %v |%s %-7s", wrapped.status, time.Since(start), GetRemoteIP(req), req.Method, req.URL.EscapedPath())
	})
}
