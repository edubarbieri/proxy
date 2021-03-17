package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

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
