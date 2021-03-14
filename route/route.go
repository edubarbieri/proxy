package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"httpproxy/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/thoas/stats"
)

type Route struct {
	Pattern      string
	BackendURI   []*url.URL
	firstHandler http.Handler
	Middlewares  []middleware.Middleware
	StatsEnabled bool
	stats        *stats.Stats
}

func NewRoute(pattern string, backendURIs []string, middlewares []middleware.Middleware, statsEnabled bool) (Route, error) {

	urls := make([]*url.URL, 0)

	for _, backendURI := range backendURIs {
		url, parseUrlError := url.Parse(backendURI)
		if parseUrlError != nil {
			msg := fmt.Sprintf("could not parser backend uri %s - %v", backendURI, parseUrlError)
			return Route{}, errors.New(msg)
		}
		urls = append(urls, url)
	}

	var statsMid *stats.Stats
	if statsEnabled {
		statsMid = stats.New()
	}

	return Route{
		Pattern:      pattern,
		BackendURI:   urls,
		StatsEnabled: statsEnabled,
		stats:        statsMid,
		firstHandler: createChainHandler(middlewares, urls, statsMid),
		Middlewares:  middlewares,
	}, nil
}

func createChainHandler(middlewares []middleware.Middleware, backendURI []*url.URL, statsMid *stats.Stats) http.Handler {
	funcHandler := new(Proxy).newProxyHandler(backendURI)
	for index := len(middlewares) - 1; index >= 0; index-- {
		currentMid := middlewares[index]
		funcHandler = currentMid.Middleware(funcHandler)
	}
	if statsMid != nil {
		return statsMid.Handler(funcHandler)
	}
	return funcHandler

}

type Proxy struct {
	proxies []*httputil.ReverseProxy
	next    uint32
}

func (proxy *Proxy) nextProxy() *httputil.ReverseProxy {
	n := atomic.AddUint32(&proxy.next, 1)
	return proxy.proxies[(int(n)-1)%len(proxy.proxies)]
}

func (proxy *Proxy) newProxyHandler(backendURI []*url.URL) http.Handler {
	for _, url := range backendURI {
		proxy.proxies = append(proxy.proxies, httputil.NewSingleHostReverseProxy(url))
	}
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		proxy.nextProxy().ServeHTTP(res, req)
	})
}

// Main method to handler request for this rout
func (route *Route) HandlerRequest(res http.ResponseWriter, req *http.Request) {
	if route.StatsEnabled && req.URL.Query().Get("op") == "stats" {
		route.writeStatsData(res)
		return
	}
	route.firstHandler.ServeHTTP(res, req)
}
func (route *Route) StatsData() *stats.Data {
	return route.stats.Data()
}

func (route *Route) writeStatsData(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(route.stats.Data())
	res.Write(b)
}
