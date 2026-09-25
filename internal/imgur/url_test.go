package imgur

import "testing"

func TestNormalizeTarget(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"cvWgXFc", "https://imgur.com/cvWgXFc"},
		{"cvWgXFc.jpg", "https://i.imgur.com/cvWgXFc.jpg"},
		{"https://i.imgur.com/cvWgXFc.jpg", "https://i.imgur.com/cvWgXFc.jpg"},
		{"https://imgur.com/abc123", "https://imgur.com/abc123"},
	}

	for _, tc := range tests {
		got, err := NormalizeTarget(tc.in)
		if err != nil {
			t.Fatalf("NormalizeTarget(%q) error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("NormalizeTarget(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeTargetRejectsNonImgur(t *testing.T) {
	_, err := NormalizeTarget("https://example.com/image.jpg")
	if err == nil {
		t.Fatal("expected error for non-imgur URL")
	}
}

func TestNormalizeTargetRejectsBadScheme(t *testing.T) {
	_, err := NormalizeTarget("ftp://i.imgur.com/foo.jpg")
	if err == nil {
		t.Fatal("expected error for ftp scheme")
	}
}

func TestNormalizeTargetStripsUserinfo(t *testing.T) {
	got, err := NormalizeTarget("https://user:pass@i.imgur.com/foo.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://i.imgur.com/foo.jpg" {
		t.Fatalf("got %q", got)
	}
}
