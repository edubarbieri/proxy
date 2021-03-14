package main

import (
	"context"
	"errors"
	"httpproxy/middleware"
	"httpproxy/route"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var routes map[string]route.Route
var redisClient *redis.Client

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
	initEnv()
	initRedisClient()
	initRoutes()
	http.HandleFunc("/", handlerRequest)
	log.Printf("Starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
	routes = map[string]route.Route{}
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
