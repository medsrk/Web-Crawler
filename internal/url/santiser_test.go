package url

import (
	"strings"
	"testing"
)

func TestSanitiseURL_URLTooLong(t *testing.T) {
	s := NewSanitiser()
	s.MaxURLLength = 10

	longURL := "http://example.com/this/is/a/very/long/url"
	_, err := s.SanitiseURL(longURL)

	if err == nil {
		t.Error("Expected error for URL too long")
	}

	if err != nil && !strings.Contains(err.Error(), "URL length") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
