package main

import (
	"os"
	"time"
)

func initPubSub() {
	go publishStats()
}

func publishStats() {
	channel := os.Getenv("STATS_CHANNEL")
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		statsMid.PublishStatsRedis(redisClient, channel)
	}
}
