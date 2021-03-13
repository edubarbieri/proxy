package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var routes map[string]string

func getRequestURI(req *http.Request) string {
	return req.RequestURI
}

func determineRoute(uri string) (string, error) {
	if backend, hasKey := routes[uri]; hasKey {
		return backend, nil
	}
	return "", errors.New("Not found backend")
}

func handlerRequest(res http.ResponseWriter, req *http.Request) {
	requestUri := getRequestURI(req)
	log.Printf("Request URI is %s", requestUri)

	backendUrl, notBackend := determineRoute(requestUri)

	if notBackend != nil {
		log.Printf("No backend for uri %s", requestUri)
		res.WriteHeader(404)
		return
	}

	log.Printf("Backend url is %s", backendUrl)

	url, parseUrlError := url.Parse(backendUrl)

	if parseUrlError != nil {
		log.Printf("Could not parser backend url %s - %v", backendUrl, parseUrlError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}

func main() {
	routes = map[string]string{
		"/backend1": "http://localhost:3000",
		"/backend2": "http://localhost:3001",
	}
	http.HandleFunc("/", handlerRequest)
	log.Printf("Starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
