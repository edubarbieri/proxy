package middleware

import (
	"net/http"
	"testing"
)

func TestGetRemoteIpWithProxy(t *testing.T) {
	expectation := "203.0.113.195"
	req := &http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"203.0.113.195"},
		},
	}
	result := GetRemoteIP(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}

func TestGetRemoteIpWithProxyMultiIps(t *testing.T) {
	expectation := "203.0.113.195"
	req := &http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"203.0.113.195, 70.41.3.18, 150.172.238.178"},
		},
	}
	result := GetRemoteIP(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}

func TestGetRemoteIpWithoutProxy(t *testing.T) {
	expectation := "203.0.113.195"
	req := &http.Request{
		RemoteAddr: "203.0.113.195:12334",
	}
	result := GetRemoteIP(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}
func TestGetRemoteIpWithoutProxyInvalidRemoteAddress(t *testing.T) {
	expectation := ""
	req := &http.Request{
		RemoteAddr: "203.0.113.195",
	}
	result := GetRemoteIP(req)

	if result != expectation {
		t.Errorf("Expected %v but got %v", expectation, result)
	}
}
