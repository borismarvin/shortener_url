package entity

import (
	"fmt"
	"net/url"
)

type URL struct {
	Scheme string
	Host   string
	Path   string
}

func NewURL(inputURL string) (*URL, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}

	newURL := &URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}
	return newURL, nil
}

// Parses URL
//
// If URL couldn't be parsed, returns nil
func ParseURL(inputURL string) *URL {
	u, err := NewURL(inputURL)
	if err != nil {
		return &URL{}
	}

	return u
}

// Validates URL
func IsValidURL(inputURL string) bool {
	if len(inputURL) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(inputURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func (u URL) String() string {
	s, err := url.JoinPath(u.Host, u.Path)
	if err != nil {
		return ""
	}

	if u.Scheme != "" {
		s = fmt.Sprintf("%s://%s", u.Scheme, s)
	}

	return s
}
