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
