package middleware

import (
	"net/http"
	"net/textproto"
	"testing"
)

func TestRateLimitRuleMatchByAll(t *testing.T) {
	rateRule := RateLimiteRule{
		ID:          "test",
		Limit:       100,
		TargetPath:  "/test",
		HeaderValue: "X-Api-Key",
	}
	req := &http.Request{
		RequestURI: "/test/service/123",
		Header: http.Header{
			textproto.CanonicalMIMEHeaderKey("X-Api-Key"): []string{"EMPhafKJFGZxnWmaJ97e3U8"},
		},
	}
	expectation := true

	result := rateRule.match(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}
func TestRateLimitRuleMatchByPath(t *testing.T) {
	rateRule := RateLimiteRule{
		ID:         "test",
		Limit:      100,
		TargetPath: "/test",
	}
	req := &http.Request{
		RequestURI: "/test/service/123",
		Header: http.Header{
			textproto.CanonicalMIMEHeaderKey("X-ApiKey"): []string{"test"},
		},
	}
	expectation := true

	result := rateRule.match(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}
func TestRateLimitRuleMatchByHeader(t *testing.T) {
	rateRule := RateLimiteRule{
		ID:          "test",
		Limit:       100,
		HeaderValue: "X-ApiKey",
	}
	req := &http.Request{
		RequestURI: "/test/service/123",
		Header: http.Header{
			textproto.CanonicalMIMEHeaderKey("X-ApiKey"): []string{"EMPhafKJFGZxnWmaJ97e3U8"},
		},
	}
	expectation := true

	result := rateRule.match(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}

// func TestRateKey(t *testing.T) {
// 	md := NewRateLimitMiddleware(&redis.Client{}, "/test", 10, false, "")
// 	req := &http.Request{
// 		Header: http.Header{
// 			"X-Forwarded-For": []string{"203.0.113.195"},
// 		},
// 	}
// 	minute := time.Now().Format("15:04")
// 	expectation := "rate_/test_" + minute

// 	result := md.rateKey(req)

// 	if result != expectation {
// 		t.Errorf("Expected %v but got %v", expectation, result)
// 	}
// }
// func TestRateKeyWithIp(t *testing.T) {
// 	md := NewRateLimitMiddleware(&redis.Client{}, "/test", 10, true, "")
// 	req := &http.Request{
// 		Header: http.Header{
// 			"X-Forwarded-For": []string{"203.0.113.195"},
// 		},
// 	}
// 	minute := time.Now().Format("15:04")
// 	expectation := "rate_/test_203.0.113.195_" + minute

// 	result := md.rateKey(req)

// 	if result != expectation {
// 		t.Errorf("Expected %v but got %v", expectation, result)
// 	}
// }

// func TestRateKeyWithHeader(t *testing.T) {
// 	md := NewRateLimitMiddleware(&redis.Client{}, "/test", 10, false, "X-Api-Key")
// 	req := &http.Request{
// 		Header: http.Header{
// 			"X-Api-Key": []string{"EMPhafKJFGZxnWmaJ97e3U8"},
// 		},
// 	}
// 	minute := time.Now().Format("15:04")
// 	expectation := "rate_/test_EMPhafKJFGZxnWmaJ97e3U8_" + minute

// 	result := md.rateKey(req)

// 	if result != expectation {
// 		t.Errorf("Expected %v but got %v", expectation, result)
// 	}
// }
