package config

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Ports     Ports      `yaml:"ports" validate:"required"`
		LogFile   string     `yaml:"logFile" validate:"required"`
		CertsFile string     `yaml:"certsFile" validate:"required"`
		LogLevel  string     `yaml:"logLevel" validate:"oneof='debug' 'info' 'warn' 'error'"`
		Endpoints []Endpoint `yaml:"endpoints" validate:"dive"`
	}

	Ports struct {
		HTTP  int `yaml:"http" validate:"required,gt=0"`
		HTTPS int `yaml:"https" validate:"gte=0"`
	}

	Endpoint struct {
		URL             URL    `yaml:"url" validate:"required"`
		HTTPS           string `yaml:"https" validate:"oneof='' 'acme' 'self'"`
		RedirectToHTTPS bool   `yaml:"redirectToHttps"`
		Enabled         *bool  `yaml:"enabled"`

		Redirect   string      `yaml:"redirect" validate:"required_without_all=Static Proxy RunCommand"`
		Static     *Static     `yaml:"static" validate:"required_without_all=Redirect Proxy RunCommand"`
		Proxy      *Proxy      `yaml:"proxy" validate:"required_without_all=Static Redirect RunCommand"`
		RunCommand *RunCommand `yaml:"runCommand" validate:"required_without_all=Static Proxy Redirect"`
	}

	Static struct {
		Dir      string `yaml:"dir" validate:"required"`
		Index    string `yaml:"index" validate:"required"`
		NotFound string `yaml:"notFound" validate:"required"`
	}

	Proxy struct {
		URL          URL    `yaml:"url"`
		RemovePrefix string `yaml:"removePrefix"`
	}

	RunCommand struct {
		Token   string   `yaml:"token" validate:"required"`
		Command []string `yaml:"command" validate:"required,min=1"`
	}
)

var defaultConfig = Config{
	Ports: Ports{
		HTTP:  80,
		HTTPS: 443,
	},
	LogFile:   "requests.log",
	CertsFile: "certs.json",
	LogLevel:  "debug",
}

var validatorInstance = validator.New()

const (
	defaultIndex    = "index.html"
	defaultNotFound = "404.html"
)

func (config *Config) LoadFromYAML(data []byte) error {
	*config = defaultConfig
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	for i, endpoint := range config.Endpoints {
		if endpoint.Static != nil {
			if endpoint.Static.Index == "" {
				endpoint.Static.Index = defaultIndex
			}
			if endpoint.Static.NotFound == "" {
				endpoint.Static.NotFound = defaultNotFound
			}
			config.Endpoints[i] = endpoint
		}
	}
	return validatorInstance.Struct(config)
}

func (config *Config) AllHosts() []string {
	hostsMap := make(map[string]struct{}, len(config.Endpoints))
	for _, endpoint := range config.Endpoints {
		if endpoint.Enabled != nil && !*endpoint.Enabled {
			continue
		}
		hostsMap[endpoint.URL.Hostname()] = struct{}{}
	}

	hosts := make([]string, 0, len(hostsMap))
	for host := range hostsMap {
		hosts = append(hosts, host)
	}
	return hosts
}

func (config *Config) AcmeHosts() []string {
	hostsMap := make(map[string]struct{}, len(config.Endpoints))
	for _, endpoint := range config.Endpoints {
		if endpoint.Enabled != nil && !*endpoint.Enabled {
			continue
		}
		if endpoint.HTTPS != "acme" {
			continue
		}
		hostsMap[endpoint.URL.Hostname()] = struct{}{}
	}

	hosts := make([]string, 0, len(hostsMap))
	for host := range hostsMap {
		hosts = append(hosts, host)
	}
	return hosts
}

func (config *Config) SelfHosts() []string {
	hostsMap := make(map[string]struct{}, len(config.Endpoints))
	for _, endpoint := range config.Endpoints {
		if endpoint.Enabled != nil && !*endpoint.Enabled {
			continue
		}
		if endpoint.HTTPS != "self" {
			continue
		}
		hostsMap[endpoint.URL.Hostname()] = struct{}{}
	}

	hosts := make([]string, 0, len(hostsMap))
	for host := range hostsMap {
		hosts = append(hosts, host)
	}
	return hosts
}

func (config *Config) HTTPS() bool {
	for _, endpoint := range config.Endpoints {
		if endpoint.HTTPS != "" {
			return true
		}
	}

	return false
}

func (config *Config) GetLogLevel() slog.Level {
	switch config.LogLevel {
	case "", "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		panic("unknown log level")
	}
}
