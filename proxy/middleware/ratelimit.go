package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitMiddleware struct {
	redisClient *redis.Client
	rules       []RateLimiteRule
}

type RateLimiteRule struct {
	ID          string
	Limit       uint64
	TargetPath  string
	HeaderValue string
	SourceIP    bool
}

func (rule *RateLimiteRule) match(req *http.Request) bool {
	if len(rule.TargetPath) > 0 && !strings.HasPrefix(req.RequestURI, rule.TargetPath) {
		return false
	}
	if len(rule.HeaderValue) > 0 {
		headerValue := req.Header.Get(rule.HeaderValue)
		if len(headerValue) <= 0 {
			return false
		}
	}
	return true
}

func (rule *RateLimiteRule) ruleKey(req *http.Request) string {
	key := "rate:" + rule.ID
	if len(rule.TargetPath) > 0 && strings.HasPrefix(req.RequestURI, rule.TargetPath) {
		key = key + ":" + rule.TargetPath
	}
	if len(rule.HeaderValue) > 0 {
		headerValue := req.Header.Get(rule.HeaderValue)
		if len(headerValue) > 0 {
			key = key + ":" + headerValue
		}
	}
	if rule.SourceIP {
		key = key + ":" + GetRemoteIP(req)
	}

	return key
}

func NewRateLimitMiddleware(redisClient *redis.Client) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient: redisClient,
	}
}

func (md *RateLimitMiddleware) UpdateRules(rules []RateLimiteRule) {
	log.Println("updating ratelimit rules")
	md.rules = rules
}

func (md *RateLimitMiddleware) findRule(req *http.Request) (*RateLimiteRule, bool) {
	for _, rule := range md.rules {
		if rule.match(req) {
			return &rule, true
		}
	}
	return nil, false
}

func (md *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rule, notFound := md.findRule(req)
		if !notFound {
			next.ServeHTTP(res, req)
			return
		}
		ok, error := md.validateRateRule(rule, req)

		if error != nil {
			internalError(error, res)
			return
		}
		if !ok {
			res.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(res, req)
	})
}

func (md *RateLimitMiddleware) validateRateRule(rule *RateLimiteRule, req *http.Request) (bool, error) {
	key := rule.ruleKey(req) + "_" + time.Now().Format("15:04")
	currentVal, getError := md.redisClient.Get(req.Context(), key).Result()
	if getError != nil && getError != redis.Nil {
		return false, fmt.Errorf("error getting redis key value %v", getError)
	}
	if getError == nil {
		//value already exists in redis
		curValInt, parseError := strconv.ParseUint(currentVal, 10, 64)
		if parseError != nil {
			return false, fmt.Errorf("error parsing current redis key value %v", parseError)
		}

		if curValInt > rule.Limit {
			return false, nil
		}
	}
	incVal, incError := md.redisClient.Incr(req.Context(), key).Result()
	if incError != nil {
		return false, fmt.Errorf("error inc current redis key value %v", incError)
	}
	if incVal == 1 {
		_, expError := md.redisClient.Expire(req.Context(), key, time.Minute).Result()
		if expError != nil {
			return false, fmt.Errorf("error set expiration redis key value %v", expError)
		}
	}

	return true, nil
}

func internalError(err error, res http.ResponseWriter) {
	log.Println(err)
	res.WriteHeader(http.StatusInternalServerError)
}
