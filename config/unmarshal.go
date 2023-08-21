package config

import "errors"

func (s *Hosts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var hosts []string
	if err := unmarshal(&hosts); err == nil {
		*s = hosts
		return nil
	}

	var host string
	if err := unmarshal(&host); err == nil {
		*s = []string{host}
		return nil
	}

	return errors.New("can not unmarshal host")
}

func (proxy *Proxy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var url string
	if err := unmarshal(&url); err == nil {
		*proxy = Proxy{URL: url}
		return nil
	}

	var pr struct {
		URL          string `yaml:"url"`
		RemovePrefix string `yaml:"removePrefix"`
		AddPrefix    string `yaml:"addPrefix"`
	}
	if err := unmarshal(&pr); err == nil {
		*proxy = pr
		return nil
	}

	return errors.New("can not unmarshal proxy")
}
