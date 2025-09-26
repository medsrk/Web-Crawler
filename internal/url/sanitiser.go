package url

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// errors
var (
	ErrInvalidURL        = errors.New("invalid URL")
	ErrUnsupportedScheme = errors.New("unsupported scheme")
	ErrURLTooLong        = errors.New("URL is too long")
	ErrEmptyURL          = errors.New("URL is empty")
)

var (
	defaultMaxLength = 2048
	defaultSchemes   = []string{"http", "https"}
)

type Sanitiser struct {
	MaxURLLength   int
	AllowedSchemes []string
	DefaultScheme  string
}

func NewSanitiser() *Sanitiser {
	return &Sanitiser{
		MaxURLLength:   defaultMaxLength,
		AllowedSchemes: defaultSchemes,
		DefaultScheme:  defaultSchemes[1],
	}
}

func (s *Sanitiser) SanitiseURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", ErrEmptyURL
	}

	if len(rawURL) > s.MaxURLLength {
		return "", fmt.Errorf("%w: URL length %d exceeds maximum of %d characters %s", ErrURLTooLong, len(rawURL), s.MaxURLLength, rawURL)
	}

	normalised := s.normaliseURL(rawURL)

	parsed, err := url.Parse(normalised)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if !s.isSchemeAllowed(parsed.Scheme) {
		return "", fmt.Errorf("%w: scheme '%s' is not allowed, allowed schemes are: %v", ErrUnsupportedScheme, parsed.Scheme, s.AllowedSchemes)
	}

	return s.cleanURL(parsed).String(), nil

}

func (s *Sanitiser) normaliseURL(rawURL string) string {
	var normalisedURL string
	normalisedURL = strings.TrimSpace(rawURL)

	if !strings.Contains(normalisedURL, "://") {
		normalisedURL = s.DefaultScheme + rawURL
	}

	return normalisedURL
}

func (s *Sanitiser) isSchemeAllowed(scheme string) bool {
	scheme = strings.ToLower(scheme)
	for _, allowed := range s.AllowedSchemes {
		if scheme == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}

func (s *Sanitiser) cleanURL(parsed *url.URL) *url.URL {
	return &url.URL{
		Scheme:   strings.ToLower(parsed.Scheme),
		Host:     strings.ToLower(parsed.Host),
		Path:     s.cleanPath(parsed.Path),
		RawQuery: parsed.RawQuery,
	}
}

func (s *Sanitiser) cleanPath(path string) string {
	if path == "" {
		return "/"
	}

	return strings.ReplaceAll(path, "//", "/")
}
