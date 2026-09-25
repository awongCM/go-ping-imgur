package imgur

import (
	"fmt"
	"net/http"
	"time"
)

// Result captures one HTTP ping attempt.
type Result struct {
	Target     string
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func (r Result) OK() bool {
	return r.Err == nil && r.StatusCode >= 200 && r.StatusCode < 400
}

// Ping issues a HEAD request to the normalized URL.
func Ping(client *http.Client, target string) Result {
	result := Result{Target: target}

	url, err := NormalizeTarget(target)
	if err != nil {
		result.Err = err
		return result
	}
	result.URL = url

	if client == nil {
		client = NewClient(15 * time.Second)
	}

	start := time.Now()
	resp, err := client.Head(url)
	result.Duration = time.Since(start)
	if err != nil {
		result.Err = fmt.Errorf("request failed: %w", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if result.StatusCode >= 400 {
		result.Err = fmt.Errorf("HTTP %d", result.StatusCode)
	}

	return result
}
