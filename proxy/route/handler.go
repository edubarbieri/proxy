package route

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/edubarbieri/proxy/config"
	"github.com/edubarbieri/proxy/middleware"
)

type EntryPoint struct {
	routes map[string]*Route
}

func (e *EntryPoint) UpdateConfig(c config.Config) {
	log.Println("updating entry point route config")
	e.routes = e.createRoutes(c.Routes)
}

func (e *EntryPoint) determineRoute(uri string) (*Route, error) {
	for pattern, route := range e.routes {
		if strings.HasPrefix(uri, pattern) {
			return route, nil
		}
	}
	return nil, errors.New("not found backend")
}
func (e *EntryPoint) HandlerRequest(res http.ResponseWriter, req *http.Request) {
	requestUri := req.RequestURI
	route, noRoute := e.determineRoute(requestUri)

	if noRoute != nil {
		log.Printf("No backend for uri %s", requestUri)
		res.WriteHeader(404)
		return
	}
	route.HandlerRequest(res, req)
}

func (e *EntryPoint) createRoutes(routesConfigs []config.RouteConfig) map[string]*Route {
	routesMap := map[string]*Route{}
	for _, routeConf := range routesConfigs {
		route, error := e.createRoute(routeConf)
		if error != nil {
			log.Printf("error creating route %s - %v", routeConf.Pattern, error)
			continue
		}
		routesMap[routeConf.Pattern] = route
	}

	return routesMap
}

func (e *EntryPoint) createRoute(routeConfig config.RouteConfig) (*Route, error) {
	return NewRoute(routeConfig.Pattern, routeConfig.Backends, []middleware.Middleware{})
}
