package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

const secretToken = "<secret>"

// Secret is a string that must not be revealed on marshaling.
type Secret string

// MarshalYAML implements the yaml.Marshaler interface for Secret.
func (s Secret) MarshalYAML() (interface{}, error) {
	if s != "" {
		return secretToken, nil
	}
	return nil, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Secret.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Secret
	return unmarshal((*plain)(s))
}

// MarshalJSON implements the json.Marshaler interface for Secret.
func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(secretToken)
}

// URL is a custom type that represents an HTTP or HTTPS URL and allows validation at configuration load time.
type URL struct {
	url.URL
}

// ParseURL constructs a new URL from url string.
func ParseURL(s string) (*URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme %q for URL", u.Scheme)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("missing host for URL")
	}

	return &URL{*u}, nil
}

// MarshalYAML implements the yaml.Marshaler interface for URL.
func (u *URL) MarshalYAML() (interface{}, error) {
	return u.URL.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for URL.
func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	urlp, err := ParseURL(s)
	if err != nil {
		return err
	}
	u.URL = urlp.URL
	return nil
}

// MarshalJSON implements the json.Marshaler interface for URL.
func (u *URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.URL.String())
}

// SecretURL is a URL that must not be revealed on marshaling.
type SecretURL URL

// Copy makes a deep-copy of the type
func (s *SecretURL) Copy() SecretURL {
	v := *s
	return v
}

var secretRE = regexp.MustCompile(`(?i)(secret|token|key|nonce|digest)`)

// MarshalYAML implements the yaml.Marshaler interface for SecretURL.
func (s SecretURL) MarshalYAML() (interface{}, error) {
	cloned := s
	qs := cloned.Query()
	keys := make([]string, 0, len(qs))
	for k := range qs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for _, k := range keys {
		keyEscaped := url.QueryEscape(k)

		if secretRE.MatchString(k) {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}

			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString("<secret>")
			continue
		}

		for _, v := range qs[k] {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	cloned.RawQuery = buf.String()
	return cloned.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for SecretURL.
func (s *SecretURL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	return unmarshal((*URL)(s))
}

// MarshalJSON implements the json.Marshaler interface for SecretURL.
func (s SecretURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(secretToken)
}
