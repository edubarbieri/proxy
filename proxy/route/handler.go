package route

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

var routes map[string]*Route

func UpdateRoutes(newRoutes map[string]*Route) {
	routes = newRoutes
}

func determineRoute(uri string) (*Route, error) {
	for pattern, route := range routes {
		if strings.HasPrefix(uri, pattern) {
			return route, nil
		}
	}
	return nil, errors.New("not found backend")
}
func HandlerRequest(res http.ResponseWriter, req *http.Request) {
	requestUri := req.RequestURI
	route, noRoute := determineRoute(requestUri)

	if noRoute != nil {
		log.Printf("No backend for uri %s", requestUri)
		res.WriteHeader(404)
		return
	}
	route.HandlerRequest(res, req)
}
