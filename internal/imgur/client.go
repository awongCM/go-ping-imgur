package imgur

import (
	"fmt"
	"net/http"
	"time"
)

// NewClient returns an HTTP client for Imgur pings with timeout and safe redirects.
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:       timeout,
		CheckRedirect: checkImgurRedirect,
	}
}

func checkImgurRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("stopped after 10 redirects")
	}
	if !IsImgurHost(req.URL.Host) {
		return fmt.Errorf("redirect blocked: non-imgur host %q", req.URL.Host)
	}
	return nil
}
