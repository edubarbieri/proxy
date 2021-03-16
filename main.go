package main

import (
	"context"
	"errors"
	"httpproxy/middleware"
	"httpproxy/route"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var mainRateLimit *middleware.RateLimitMiddleware
var routes map[string]*route.Route
var redisClient *redis.Client

func main() {
	initEnv()
	initRedisClient()
	mainRateLimit = middleware.NewRateLimitMiddleware(redisClient)
	//Log -> Stats -> RateLimit -> handlerRequest
	finalHandler := http.HandlerFunc(handlerRequest)
	firstFnc := mainRateLimit.Middleware(finalHandler)
	initRoutes()

	if os.Getenv("DEBUG") == "true" {
		logMid := middleware.LogMiddleware{}
		firstFnc = logMid.Middleware(firstFnc)
	}

	http.Handle("/", firstFnc)
	log.Printf("Starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerRequest(res http.ResponseWriter, req *http.Request) {
	requestUri := req.RequestURI
	route, noRoute := determineRoute(requestUri)

	if noRoute != nil {
		log.Printf("No backend for uri %s", requestUri)
		res.WriteHeader(404)
		return
	}
	route.HandlerRequest(res, req)
}
func determineRoute(uri string) (*route.Route, error) {
	for pattern, route := range routes {
		if strings.HasPrefix(uri, pattern) {
			return route, nil
		}
	}
	return nil, errors.New("not found backend")
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func initRedisClient() {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	log.Printf("initializing redis client for server %s...\n", redisAddress)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("error connecting to redis server: %v", err)
	}
}

func initRoutes() {
	log.Println("initializing routes configurations...")
	routes = map[string]*route.Route{}
	mids := []middleware.Middleware{}
	mids2 := []middleware.Middleware{}
	route1, _ := route.NewRoute("/backend1", []string{"http://localhost:3000"}, mids)
	route2, _ := route.NewRoute("/backend2", []string{"http://localhost:3001", "http://localhost:3002", "http://localhost:3003", "http://localhost:3004"}, mids2)

	routes[route1.Pattern] = route1
	routes[route2.Pattern] = route2

	rule := middleware.RateLimiteRule{
		ID:          "1",
		Limit:       10,
		TargetPath:  "/backend2",
		SourceIP:    true,
		HeaderValue: "api-key",
	}

	mainRateLimit.UpdateRules([]middleware.RateLimiteRule{rule})
}
