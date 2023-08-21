package config

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		LogFile   string    `yaml:"logFile"`
		CertsFile string    `yaml:"certsFile"`
		LogLevel  string    `yaml:"logLevel"`
		Services  []Service `yaml:"services"`
	}

	Service struct {
		Hosts    Hosts  `yaml:"hosts"`
		Path     string `yaml:"path"`
		Disabled bool   `yaml:"disabled"`
		Page404  string `yaml:"page404"`

		Redirect   *string     `yaml:"redirect"`
		Static     *string     `yaml:"static"`
		Proxy      *Proxy      `yaml:"proxy"`
		RunCommand *RunCommand `yaml:"runCommand"`
	}

	Hosts []string

	Proxy struct {
		URL          string `yaml:"url"`
		RemovePrefix string `yaml:"removePrefix"`
		AddPrefix    string `yaml:"addPrefix"`
	}

	RunCommand struct {
		Token   string   `yaml:"token"`
		Command []string `yaml:"command"`
	}
)

func (config *Config) LoadFromYAML(data []byte) error {
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return config.init()
}

func (config *Config) init() error {
	config.LogLevel = strings.ToLower(config.LogLevel)
	if config.LogLevel == "" {
		config.LogLevel = "debug"
	}
	if config.LogFile == "" {
		config.LogFile = "~/mws.log"
	}
	if config.CertsFile == "" {
		config.LogFile = "~/certs.json"
	}

	for i := range config.Services {
		if err := config.Services[i].init(); err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) init() error {
	count := 0

	if service.Redirect != nil {
		count++
	}
	if service.Static != nil {
		count++
	}
	if service.Proxy != nil {
		count++
	}
	if service.RunCommand != nil {
		count++
	}

	if count == 0 {
		return errors.New("no actions in a service")
	}
	return nil
}

func (config Config) CollectMapByHost() map[string][]Service {
	services := make(map[string][]Service, len(config.Services))
	for _, service := range config.Services {
		for _, host := range service.Hosts {
			services[host] = append(services[host], service)
		}
	}
	return services
}
