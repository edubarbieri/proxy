package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRateKey(t *testing.T) {
	md := NewRateLimitMiddleware(&redis.Client{}, "/test", 10, false, "")
	req := &http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"203.0.113.195"},
		},
	}
	minute := time.Now().Format("15:04")
	expectation := "rate_/test_" + minute

	result := md.rateKey(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}
func TestRateKeyWithIp(t *testing.T) {
	md := NewRateLimitMiddleware(&redis.Client{}, "/test", 10, true, "")
	req := &http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"203.0.113.195"},
		},
	}
	minute := time.Now().Format("15:04")
	expectation := "rate_/test_203.0.113.195_" + minute

	result := md.rateKey(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}

func TestRateKeyWithHeader(t *testing.T) {
	md := NewRateLimitMiddleware(&redis.Client{}, "/test", 10, false, "X-Api-Key")
	req := &http.Request{
		Header: http.Header{
			"X-Api-Key": []string{"EMPhafKJFGZxnWmaJ97e3U8"},
		},
	}
	minute := time.Now().Format("15:04")
	expectation := "rate_/test_EMPhafKJFGZxnWmaJ97e3U8_" + minute

	result := md.rateKey(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}
