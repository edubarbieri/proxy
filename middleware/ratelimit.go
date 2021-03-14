package middleware

import (
	"log"
	"net/http"
)

type RateLimitMiddleware struct {
	Name string
}

func (md *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("From RateLimitMiddleware")
		next.ServeHTTP(res, req)
	})
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{}
}
