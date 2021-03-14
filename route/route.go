package route

import (
	"errors"
	"fmt"
	"httpproxy/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Route struct {
	Pattern      string
	BackendURI   *url.URL
	firstHandler http.Handler
}

func NewRoute(pattern string, backendURI string, middlewares []middleware.Middleware) (Route, error) {
	url, parseUrlError := url.Parse(backendURI)

	if parseUrlError != nil {
		msg := fmt.Sprintf("could not parser backend uri %s - %v", backendURI, parseUrlError)
		return Route{}, errors.New(msg)
	}
	return Route{
		Pattern:      pattern,
		BackendURI:   url,
		firstHandler: createChainHandler(middlewares, url),
	}, nil
}

// Main method to handler request for this rout
func (route *Route) HandlerRequest(res http.ResponseWriter, req *http.Request) {
	route.firstHandler.ServeHTTP(res, req)
}

func createChainHandler(middlewares []middleware.Middleware, backendURI *url.URL) http.Handler {
	funcHandler := newProxyHandler(backendURI)
	if len(middlewares) == 0 {
		return funcHandler
	}

	for index := len(middlewares) - 1; index >= 0; index-- {
		currentMid := middlewares[index]
		funcHandler = currentMid.Middleware(funcHandler)
	}

	return funcHandler

}

func newProxyHandler(backendURI *url.URL) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(backendURI)
		proxy.ServeHTTP(res, req)
	})
}
