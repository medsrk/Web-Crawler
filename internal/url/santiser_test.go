package url

import (
	"errors"
	"strings"
	"testing"
)

func TestSanitiseURL_EmptyURL(t *testing.T) {
	s := NewSanitiser()

	_, err := s.SanitiseURL("")
	if !errors.Is(err, ErrEmptyURL) {
		t.Errorf("Expected ErrEmptyURL, got %v", err)
	}
}

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

func TestSanitiseURL_ValidCases(t *testing.T) {
	s := NewSanitiser()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid HTTPS URL",
			input:    "https://example.com/path",
			expected: "https://example.com/path",
		},
		{
			name:     "valid HTTP URL",
			input:    "http://example.com/path",
			expected: "http://example.com/path",
		},
		{
			name:     "URL without scheme gets default",
			input:    "example.com/path",
			expected: "https://example.com/path",
		},
		{
			name:     "URL with whitespace gets trimmed",
			input:    "  https://example.com/path  ",
			expected: "https://example.com/path",
		},
		{
			name:     "uppercase scheme and host get lowercased",
			input:    "HTTPS://EXAMPLE.COM/Path",
			expected: "https://example.com/Path",
		},
		{
			name:     "empty path gets slash",
			input:    "https://example.com",
			expected: "https://example.com/",
		},
		{
			name:     "double slashes in path get cleaned",
			input:    "https://example.com//path//to//resource",
			expected: "https://example.com/path/to/resource",
		},
		{
			name:     "URL with query parameters",
			input:    "https://example.com/path?param=value&other=test",
			expected: "https://example.com/path?param=value&other=test",
		},
		{
			name:     "case insensitive scheme matching",
			input:    "HTTP://example.com",
			expected: "http://example.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.SanitiseURL(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSanitiseURL_ErrorCases(t *testing.T) {
	s := NewSanitiser()

	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "unsupported scheme",
			input:         "ftp://example.com",
			expectedError: "scheme 'ftp' is not allowed",
		},
		{
			name:          "invalid URL format",
			input:         "://invalid-url",
			expectedError: "invalid URL",
		},
		{
			name:          "malformed URL",
			input:         "http://[::1:80",
			expectedError: "invalid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.SanitiseURL(tt.input)
			if err == nil {
				t.Errorf("Expected error for %s", tt.name)
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', got %v", tt.expectedError, err)
			}
		})
	}
}

func TestIsSchemeAllowed(t *testing.T) {
	s := NewSanitiser()

	tests := []struct {
		name     string
		scheme   string
		expected bool
	}{
		{
			name:     "http allowed",
			scheme:   "http",
			expected: true,
		},
		{
			name:     "https allowed",
			scheme:   "https",
			expected: true,
		},
		{
			name:     "HTTP uppercase allowed",
			scheme:   "HTTP",
			expected: true,
		},
		{
			name:     "HTTPS uppercase allowed",
			scheme:   "HTTPS",
			expected: true,
		},
		{
			name:     "ftp not allowed",
			scheme:   "ftp",
			expected: false,
		},
		{
			name:     "file not allowed",
			scheme:   "file",
			expected: false,
		},
		{
			name:     "custom not allowed",
			scheme:   "custom",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.isSchemeAllowed(tt.scheme)
			if result != tt.expected {
				t.Errorf("Expected %t for scheme '%s', got %t", tt.expected, tt.scheme, result)
			}
		})
	}
}

func TestCleanPath(t *testing.T) {
	s := NewSanitiser()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty path returns slash",
			input:    "",
			expected: "/",
		},
		{
			name:     "normal path unchanged",
			input:    "/path/to/resource",
			expected: "/path/to/resource",
		},
		{
			name:     "double slashes cleaned",
			input:    "//path//to//resource",
			expected: "/path/to/resource",
		},
		{
			name:     "multiple consecutive slashes cleaned",
			input:    "///path///to///resource",
			expected: "/path/to/resource",
		},
		{
			name:     "single slash unchanged",
			input:    "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.cleanPath(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCustomSanitiser_CustomSchemes(t *testing.T) {
	s := &Sanitiser{
		MaxURLLength:   2048,
		AllowedSchemes: []string{"custom", "special"},
		DefaultScheme:  "custom://",
	}

	result, err := s.SanitiseURL("custom://example.com/path")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "custom://example.com/path"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
