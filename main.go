package main

import (
	"errors"
	"httpproxy/middleware"
	"httpproxy/route"
	"log"
	"net/http"
)

var routes map[string]route.Route

func getRequestURI(req *http.Request) string {
	return req.RequestURI
}

func determineRoute(uri string) (route.Route, error) {
	if backend, hasKey := routes[uri]; hasKey {
		return backend, nil
	}
	return route.Route{}, errors.New("not found backend")
}

func handlerRequest(res http.ResponseWriter, req *http.Request) {
	requestUri := getRequestURI(req)
	route, noRoute := determineRoute(requestUri)

	if noRoute != nil {
		log.Printf("No backend for uri %s", requestUri)
		res.WriteHeader(404)
		return
	}
	route.HandlerRequest(res, req)
}

func main() {
	routes = map[string]route.Route{}
	initRoutes()
	http.HandleFunc("/", handlerRequest)
	log.Printf("Starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initRoutes() {
	mids := []middleware.Middleware{
		middleware.NewStatsMiddleware(),
		middleware.NewRateLimitMiddleware(),
		middleware.NewLogMiddleware("Backend 1"),
	}
	mids2 := []middleware.Middleware{
		middleware.NewLogMiddleware("Backend 2"),
	}
	route1, _ := route.NewRoute("/backend1", "http://localhost:3000", mids)
	route2, _ := route.NewRoute("/backend2", "http://localhost:3001", mids2)

	routes[route1.Pattern] = route1
	routes[route2.Pattern] = route2
}
