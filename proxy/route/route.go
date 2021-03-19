package route

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/edubarbieri/proxy/middleware"
)

type Route struct {
	Pattern      string
	BackendURIs  []*url.URL
	firstHandler http.Handler
	Middlewares  []middleware.Middleware
	proxies      []*httputil.ReverseProxy
	next         uint64
	isInited     bool
	sync.Mutex
}

func NewRoute(pattern string, backendURIs []string, middlewares []middleware.Middleware) (*Route, error) {
	urls := make([]*url.URL, 0)
	for _, backendURI := range backendURIs {
		url, parseUrlError := url.Parse(backendURI)
		if parseUrlError != nil {
			msg := fmt.Sprintf("could not parser backend uri %s - %v", backendURI, parseUrlError)
			return nil, errors.New(msg)
		}
		urls = append(urls, url)
	}
	return &Route{
		Pattern:     pattern,
		BackendURIs: urls,
		Middlewares: middlewares,
	}, nil
}

// Main method to handler request for this rout
func (r *Route) HandlerRequest(res http.ResponseWriter, req *http.Request) {
	if !r.isInited {
		r.Init()
	}
	r.firstHandler.ServeHTTP(res, req)
}

func (r *Route) Init() {
	r.Lock()
	defer r.Unlock()
	if r.isInited {
		return
	}
	log.Printf("initializing route %v", r.Pattern)
	funcHandler := r.initProxy()
	for index := len(r.Middlewares) - 1; index >= 0; index-- {
		currentMid := r.Middlewares[index]
		funcHandler = currentMid.Middleware(funcHandler)
	}
	r.firstHandler = funcHandler
	r.isInited = true
}

func (r *Route) initProxy() http.Handler {
	r.next = 0
	r.proxies = make([]*httputil.ReverseProxy, 0)
	for _, url := range r.BackendURIs {
		r.proxies = append(r.proxies, httputil.NewSingleHostReverseProxy(url))
	}
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		r.nextProxy().ServeHTTP(res, req)
	})
}

func (r *Route) nextProxy() *httputil.ReverseProxy {
	n := atomic.AddUint64(&r.next, 1)
	return r.proxies[(int(n)-1)%len(r.proxies)]
}
