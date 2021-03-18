package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/edubarbieri/proxy/config"
	"github.com/edubarbieri/proxy/middleware"
	"github.com/edubarbieri/proxy/route"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var redisClient *redis.Client

func main() {
	godotenv.Load()
	initRedisClient()
	//Config manager
	configManager := &config.ConfigManager{}
	rateLimitMid := middleware.NewRateLimitMiddleware(redisClient)
	configManager.AddUpdateObserver(rateLimitMid)
	statsMid := middleware.NewStatsMiddleware()
	adminMind := NewAdminMiddleware(statsMid, configManager)
	//Entry point with all routes
	entryPoint := &route.EntryPoint{}
	configManager.AddUpdateObserver(entryPoint)
	//Log -> Stats -> RateLimit -> Admin -> EntryPint
	routeHandler := http.HandlerFunc(entryPoint.HandlerRequest)
	firstFnc := statsMid.Middleware(rateLimitMid.Middleware(adminMind.Middleware(routeHandler)))
	if os.Getenv("DEBUG") == "true" {
		logMid := middleware.LogMiddleware{}
		firstFnc = logMid.Middleware(firstFnc)
	}
	http.Handle("/", firstFnc)
	configManager.Init()
	log.Printf("starting proxy in port %v", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
