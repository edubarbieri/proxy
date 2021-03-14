package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitMiddleware struct {
	redisClient   *redis.Client
	route         string
	limit         uint64
	limitByIP     bool
	limitByHeader string
}

func NewRateLimitMiddleware(redisClient *redis.Client, route string, limit uint64, limitByIP bool, limitByHeader string) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient:   redisClient,
		limit:         limit,
		route:         route,
		limitByIP:     limitByIP,
		limitByHeader: limitByHeader,
	}
}

func (md *RateLimitMiddleware) rateKey(req *http.Request) string {
	minute := time.Now().Format("15:04")
	key := "rate_" + md.route
	if md.limitByIP {
		key += "_" + GetRemoteIP(req)
	}
	if len(md.limitByHeader) > 0 {
		key += "_" + req.Header.Get(md.limitByHeader)
	}
	return key + "_" + minute
}

func internalError(msg string, err error, res http.ResponseWriter) {
	log.Printf(msg, err)
	res.WriteHeader(http.StatusInternalServerError)
}

func (md *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		key := md.rateKey(req)
		currentVal, getError := md.redisClient.Get(req.Context(), key).Result()
		if getError != nil && getError != redis.Nil {
			internalError("error getting redis key value %v", getError, res)
			return
		}
		if getError == nil {
			//value already exists in redis
			curValInt, parseError := strconv.ParseUint(currentVal, 10, 64)
			if parseError != nil {
				internalError("error parsing current redis key value %v", parseError, res)
				return
			}

			if curValInt > md.limit {
				log.Printf("exceed limit key %v", key)
				res.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}
		incVal, incError := md.redisClient.Incr(req.Context(), key).Result()
		if incError != nil {
			internalError("error inc current redis key value %v", incError, res)
			return
		}
		if incVal == 1 {
			_, expError := md.redisClient.Expire(req.Context(), key, time.Minute).Result()
			if expError != nil {
				internalError("error set expiration redis key value %v", incError, res)
				return
			}
		}
		next.ServeHTTP(res, req)
	})
}
