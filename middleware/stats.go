package middleware

import (
	"log"
	"net/http"
)

type StatsMiddleware struct {
	Name string
}

func (md *StatsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("From StatsMiddleware")
		next.ServeHTTP(res, req)
	})
}
func NewStatsMiddleware() *StatsMiddleware {
	return &StatsMiddleware{}
}
