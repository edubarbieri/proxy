package main

import (
	"httpproxy/middleware"
	"httpproxy/route"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var mainStatsMid *middleware.StatsMiddleware
var mainRateLimit *middleware.RateLimitMiddleware
var redisClient *redis.Client

func main() {
	initEnv()
	initRedisClient()
	initRoutes()
	initPubSub()

	http.Handle("/", initHttpHandlers())
	log.Printf("starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func initHttpHandlers() http.Handler {
	mainRateLimit = middleware.NewRateLimitMiddleware(redisClient)
	mainStatsMid = middleware.NewStatsMiddleware()
	//Log -> Stats -> RateLimit -> handleHandlerRequestrRequest
	routeHandler := http.HandlerFunc(route.HandlerRequest)
	firstFnc := mainStatsMid.Middleware(mainRateLimit.Middleware(routeHandler))
	if os.Getenv("DEBUG") == "true" {
		logMid := middleware.LogMiddleware{}
		firstFnc = logMid.Middleware(firstFnc)
	}
	return firstFnc
}

func initRoutes() {
	config := ReadConfigJson()
	log.Println("initializing routes configurations...")
	route.UpdateRoutes(CreateRoutes(config.Routes))
	log.Println("initializing rate limit configurations...")
	mainRateLimit.UpdateRules(CreateRateLimiteRules(config.Limits))
}
