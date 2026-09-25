package imgur

import (
	"fmt"
	"net/url"
	"strings"
)

// NormalizeTarget turns user input into a pingable Imgur URL.
// Accepts full URLs or bare image IDs (optionally with an extension).
func NormalizeTarget(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty target")
	}

	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("invalid URL: %w", err)
		}
		switch strings.ToLower(parsed.Scheme) {
		case "http", "https":
		default:
			return "", fmt.Errorf("unsupported URL scheme: %s", parsed.Scheme)
		}
		if !IsImgurHost(parsed.Host) {
			return "", fmt.Errorf("not an imgur URL: %s", raw)
		}
		parsed.User = nil
		return parsed.String(), nil
	}

	id := strings.TrimPrefix(raw, "/")
	if strings.Contains(id, "/") {
		return "", fmt.Errorf("invalid image id: %s", raw)
	}

	if strings.Contains(id, ".") {
		return fmt.Sprintf("https://i.imgur.com/%s", id), nil
	}

	// Gallery pages respond reliably for bare IDs without a known extension.
	return fmt.Sprintf("https://imgur.com/%s", id), nil
}

// IsImgurHost reports whether host is an Imgur-owned hostname.
func IsImgurHost(host string) bool {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	return host == "imgur.com" ||
		host == "www.imgur.com" ||
		host == "i.imgur.com" ||
		strings.HasSuffix(host, ".imgur.com")
}
