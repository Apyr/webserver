package internal

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	CertsDir        string `yaml:"certsDir"`
	HttpPort        int    `yaml:"httpPort"`
	HttpsPort       int    `yaml:"httpsPort"`
	RedirectToHttps *bool  `yaml:"redirectToHttps"`
	Services        yaml.MapSlice
}

func (cfg *Config) check(dir string) error {
	certs := cfg.CertsDir
	if cfg.CertsDir == "" {
		certs = "certs"
	}
	certs, err := absPath(certs, dir)
	if err != nil {
		return err
	}
	cfg.CertsDir = certs
	if cfg.RedirectToHttps == nil {
		t := true
		cfg.RedirectToHttps = &t
	}
	return nil
}

type Service struct {
	Enabled   *bool
	Endpoints []Endpoint
}

func (s Service) GetEnabled() bool {
	if s.Enabled == nil {
		return true
	}
	return *s.Enabled
}

type ServiceEx struct {
	Service
	Name     string
	FileName string
}

type ConfigEx struct {
	Config
	ExtraServices []ServiceEx
}
