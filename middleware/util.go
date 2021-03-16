package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func GetRemoteIP(req *http.Request) string {
	xForward := req.Header.Get("X-Forwarded-For")
	if len(xForward) > 0 {
		//Header can have multiple ips https://en.wikipedia.org/wiki/X-Forwarded-For#Format
		ips := strings.Split(xForward, ", ")
		if len(ips) > 0 {
			//return first client ip
			return ips[0]
		}
	}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("Could not determine ip for address %v", req.RemoteAddr)
		return ""
	}
	return ip
}

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
