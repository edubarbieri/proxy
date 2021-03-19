package middleware

import (
	"net/http"
	"net/textproto"
	"testing"
)

func TestRateLimitRuleMatchByAll(t *testing.T) {
	rateRule := RateLimitRule{
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
	rateRule := RateLimitRule{
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
	rateRule := RateLimitRule{
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
