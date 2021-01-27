package config

/*

import (
	"gopkg.in/yaml.v2"
)

type ListenStruct struct {
	Scheme   string
	Port     *int
	Autocert *bool
}

func (listen *ListenStruct) SetDefault() {
	if listen.Scheme == "" {
		listen.Scheme = "http"
	}
	if listen.Port == nil {
		if listen.Scheme == "http" {
			port := 80
			listen.Port = &port
		} else if listen.Scheme == "https" {
			port := 443
			listen.Port = &port
		}
	}
}

type ConfigStruct struct {
	Listen      []ListenStruct
	CertsDir    string
	ServicesDir string
}

func (cfg *ConfigStruct) SetDefault() {
	for _, l := range cfg.Listen {
		l.SetDefault()
	}
}

func ParseConfigStruct(text string) (ConfigStruct, error) {
	cfg := ConfigStruct{}
	err := yaml.Unmarshal([]byte(text), &cfg)
	return cfg, err
}

*/
