package config

import (
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

type URL struct {
	*url.URL
}

func (val *URL) UnmarshalYAML(unmarshal func(any) error) error {
	var u string
	if err := unmarshal(&u); err != nil {
		return err
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}
	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}
	val.URL = parsed
	return nil
}

func (u URL) MarshalYAML() (any, error) {
	return u.String(), nil
}

var (
	_ yaml.Marshaler   = URL{}
	_ yaml.Unmarshaler = (*URL)(nil)
)
