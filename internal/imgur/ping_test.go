package imgur

import (
	"net/http"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestPingSuccess(t *testing.T) {
	client := NewClient(time.Second)
	client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodHead {
			t.Fatalf("expected HEAD, got %s", req.Method)
		}
		if req.URL.Host != "i.imgur.com" {
			t.Fatalf("unexpected host %s", req.URL.Host)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       http.NoBody,
			Request:    req,
		}, nil
	})

	result := Ping(client, "cvWgXFc.jpg")
	if !result.OK() {
		t.Fatalf("expected OK, got %#v", result)
	}
}

func TestPingHTTPError(t *testing.T) {
	client := NewClient(time.Second)
	client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       http.NoBody,
			Request:    req,
		}, nil
	})

	result := Ping(client, "missing.jpg")
	if result.OK() {
		t.Fatal("expected failure for 404")
	}
}

func TestPingBlocksOffHostRedirect(t *testing.T) {
	client := NewClient(time.Second)
	client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     http.Header{"Location": {"https://example.com/taken"}},
			Body:       http.NoBody,
			Request:    req,
		}, nil
	})

	result := Ping(client, "cvWgXFc.jpg")
	if result.OK() {
		t.Fatal("expected blocked redirect to fail")
	}
}

func TestPingNormalizeError(t *testing.T) {
	result := Ping(NewClient(time.Second), "https://example.com/x.jpg")
	if result.OK() || result.Err == nil {
		t.Fatalf("expected normalize error, got %#v", result)
	}
}
