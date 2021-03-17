package main

import (
	"log"
	"net/http"
	"os"
	"proxy/middleware"
	"proxy/route"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var statsMid *middleware.StatsMiddleware
var rateLimitMid *middleware.RateLimitMiddleware
var redisClient *redis.Client

func main() {
	initEnv()
	initRedisClient()
	http.Handle("/", initHttpHandlers())
	initConfig()
	initPubSub()

	log.Printf("starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func initConfig() {
	config := ReadConfigJson()
	updateConfig(config)
}
func updateConfig(config Config) {
	log.Println("initializing routes configurations...")
	route.UpdateRoutes(CreateRoutes(config.Routes))
	log.Println("initializing rate limit configurations...")
	rateLimitMid.UpdateRules(CreateRateLimiteRules(config.Limits))
}

func initHttpHandlers() http.Handler {
	rateLimitMid = middleware.NewRateLimitMiddleware(redisClient)
	statsMid = middleware.NewStatsMiddleware()
	adminMind := NewAdminMiddleware(statsMid)
	//Log -> Stats -> RateLimit -> Admin -> handleHandlerRequestrRequest
	routeHandler := http.HandlerFunc(route.HandlerRequest)
	firstFnc := statsMid.Middleware(rateLimitMid.Middleware(adminMind.Middleware(routeHandler)))
	if os.Getenv("DEBUG") == "true" {
		logMid := middleware.LogMiddleware{}
		firstFnc = logMid.Middleware(firstFnc)
	}
	return firstFnc
}
