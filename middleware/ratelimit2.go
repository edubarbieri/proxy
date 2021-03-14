package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

type RateLimit2Middleware struct {
	limiter       *redis_rate.Limiter
	route         string
	limit         uint64
	limitByIP     bool
	limitByHeader string
}

func NewRateLimit2Middleware(redisClient *redis.Client, route string, limit uint64, limitByIP bool, limitByHeader string) *RateLimit2Middleware {
	return &RateLimit2Middleware{
		limiter:       redis_rate.NewLimiter(redisClient),
		limit:         limit,
		route:         route,
		limitByIP:     limitByIP,
		limitByHeader: limitByHeader,
	}
}

func (md *RateLimit2Middleware) rateKey(req *http.Request) string {
	currentTime := time.Now()
	minute := currentTime.Format("15:04")

	key := md.route
	if md.limitByIP {
		key += "_" + GetRemoteIP(req)
	}
	if len(md.limitByHeader) > 0 {
		key += "_" + req.Header.Get(md.limitByHeader)
	}
	return key + "_" + minute
}

func (md *RateLimit2Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		key := md.rateKey(req)
		allow, allowError := md.limiter.Allow(req.Context(), key, redis_rate.PerMinute(int(md.limit)))
		if allowError != nil {
			log.Printf("error checking rate limit %v", allowError)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		h := res.Header()
		h.Set("RateLimit-Remaining", strconv.Itoa(allow.Remaining))

		if allow.Allowed == 0 {
			// We are rate limited.

			seconds := int(allow.RetryAfter / time.Second)
			h.Set("RateLimit-RetryAfter", strconv.Itoa(seconds))

			log.Printf("exceed limit key %v", key)
			res.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(res, req)
	})
}
