package middleware

import (
	"log"
	"net/http"
)

type LogMiddleware struct {
	Name string
}

func (md *LogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Request URI is %s name is %s", req.RequestURI, md.Name)
		next.ServeHTTP(res, req)
	})
}

func NewLogMiddleware(name string) *LogMiddleware {
	return &LogMiddleware{Name: name}
}
