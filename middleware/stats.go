package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type StatsMiddleware struct {
	totalRequests        int
	totalResponseTime    time.Time
	totalRequestByStatus map[string]int
	totalRequestByPath   map[string]int
	sync.Mutex
}

func (md *StatsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(res)
		next.ServeHTTP(wrapped, req)

		statusCode := fmt.Sprintf("%d", wrapped.Status())
		responseTime := time.Since(start)
		md.Lock()
		defer md.Unlock()
		md.totalRequests++
		md.totalResponseTime = md.totalResponseTime.Add(responseTime)
		md.totalRequestByStatus[statusCode]++
		md.totalRequestByPath[req.URL.EscapedPath()]++
	})
}

func (md *StatsMiddleware) PublishStatsRedis(redisClient *redis.Client, channelName string) {
	data := md.GenerateStats()
	json, errJson := json.Marshal(data)
	if errJson != nil {
		log.Printf("error generating json stats %v", errJson)
		return
	}
	err := redisClient.Publish(context.Background(), channelName, string(json)).Err()
	if err != nil {
		log.Printf("error publishing stats in redis %v", err)
	}
}

func (md *StatsMiddleware) GenerateStats() Data {
	md.Lock()
	defer md.Unlock()
	now := time.Now()
	totalResponseTime := md.totalResponseTime.Sub(time.Time{})
	averageResponseTime := time.Duration(0)

	if md.totalRequests > 0 {
		avgNs := int64(totalResponseTime) / int64(md.totalRequests)
		averageResponseTime = time.Duration(avgNs)
	}
	hostName, _ := os.Hostname()
	r := Data{
		Pid:                    os.Getpid(),
		Hostname:               hostName,
		Time:                   now.String(),
		TimeUnix:               now.Unix(),
		TotalRequests:          md.totalRequests,
		TotalResponseTime:      totalResponseTime.String(),
		TotalResponseTimeSec:   totalResponseTime.Seconds(),
		AverageResponseTime:    averageResponseTime.String(),
		AverageResponseTimeSec: averageResponseTime.Seconds(),
		TotalRequestByStatus:   md.totalRequestByStatus,
		TotalRequestByPath:     md.totalRequestByPath,
	}

	md.totalRequests = 0
	md.totalResponseTime = time.Time{}
	md.totalRequestByStatus = map[string]int{}
	md.totalRequestByPath = map[string]int{}
	return r
}

func NewStatsMiddleware() *StatsMiddleware {
	return &StatsMiddleware{
		totalResponseTime:    time.Time{},
		totalRequestByStatus: map[string]int{},
		totalRequestByPath:   map[string]int{},
	}
}

type Data struct {
	Pid                    int            `json:"pid"`
	Hostname               string         `json:"hostname"`
	Time                   string         `json:"time"`
	TimeUnix               int64          `json:"unixtime"`
	TotalRequests          int            `json:"totalRequest"`
	TotalResponseTime      string         `json:"totalResponseTime"`
	TotalResponseTimeSec   float64        `json:"totalResponseTimeSec"`
	AverageResponseTime    string         `json:"averageResponseTime"`
	AverageResponseTimeSec float64        `json:"averageResponseTimeSec"`
	TotalRequestByStatus   map[string]int `json:"totalRequestByStatus"`
	TotalRequestByPath     map[string]int `json:"totalRequestByPath"`
}
