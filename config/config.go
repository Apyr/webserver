package config

import (
	"errors"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		LogFile   string     `yaml:"logFile"`
		CertsFile string     `yaml:"certsFile"`
		LogLevel  string     `yaml:"logLevel"`
		Endpoints []Endpoint `yaml:"endpoints"`
	}

	Endpoint struct {
		URL             URL   `yaml:"url"`
		HTTPS           *bool `yaml:"https"`
		RedirectToHTTPS *bool `yaml:"redirectToHttps"`
		Disabled        bool  `yaml:"disabled"`

		Redirect   string      `yaml:"redirect"`
		Static     *Static     `yaml:"static"`
		Proxy      *Proxy      `yaml:"proxy"`
		RunCommand *RunCommand `yaml:"runCommand"`
	}

	Static struct {
		Dir     string `yaml:"dir"`
		Index   string `yaml:"index"`
		Page404 string `yaml:"page404"`
	}

	Proxy struct {
		URL          URL    `yaml:"url"`
		RemovePrefix string `yaml:"removePrefix"`
	}

	RunCommand struct {
		Token   string   `yaml:"token"`
		Command []string `yaml:"command"`
	}

	URL struct {
		*url.URL
	}
)

const (
	defaultLogLevel  = "debug"
	defaultLogFile   = "requests.log"
	defaultCertsFile = "certs.json"
	defaultPage404   = "404.html"
	defaultIndex     = "index.html"
)

func (config *Config) LoadFromYAML(data []byte) error {
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return config.init()
}

func (config *Config) Hosts() []string {
	hostsMap := make(map[string]struct{}, len(config.Endpoints))
	for _, endpoint := range config.Endpoints {
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
		if *endpoint.HTTPS {
			return true
		}
	}

	return false
}

func (config *Config) init() error {
	config.LogLevel = strings.ToLower(config.LogLevel)
	if config.LogLevel == "" {
		config.LogLevel = defaultLogLevel
	}
	if config.LogFile == "" {
		config.LogFile = defaultLogFile
	}
	if config.CertsFile == "" {
		config.CertsFile = defaultCertsFile
	}

	for i := range config.Endpoints {
		if err := config.Endpoints[i].init(); err != nil {
			return err
		}
	}

	return nil
}

func (service *Endpoint) init() error {
	if service.HTTPS == nil {
		val := true
		service.HTTPS = &val
	}
	if service.RedirectToHTTPS == nil {
		val := *service.HTTPS
		service.RedirectToHTTPS = &val
	}

	return service.initActions()
}

func (service *Endpoint) initActions() error {
	count := 0

	if service.Redirect != "" {
		count++
	}
	if service.Static != nil {
		count++
		service.Static.init()
	}
	if service.Proxy != nil {
		count++
	}
	if service.RunCommand != nil {
		count++
		if err := service.RunCommand.init(); err != nil {
			return err
		}
	}

	if count == 0 {
		return errors.New("no actions in a service")
	}
	if count > 1 {
		return errors.New("too many actions in a service")
	}
	return nil
}

func (static *Static) init() {
	if static.Page404 == "" {
		static.Page404 = defaultPage404
	}
	if static.Index == "" {
		static.Index = defaultIndex
	}
}

func (run *RunCommand) init() error {
	if run.Token == "" {
		return errors.New("token must not be empty in runCommand")
	}
	if len(run.Command) == 0 {
		return errors.New("command must not be empty in runCommand")
	}
	return nil
}

func (proxy *Proxy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var url URL
	if err := unmarshal(&url); err == nil {
		*proxy = Proxy{URL: url}
		return nil
	}

	var pr struct {
		URL          URL    `yaml:"url"`
		RemovePrefix string `yaml:"removePrefix"`
	}
	if err := unmarshal(&pr); err == nil {
		*proxy = pr
		return nil
	}

	return errors.New("can not unmarshal proxy")
}

func (val *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var u string
	if err := unmarshal(&u); err != nil {
		return err
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}

	*val = URL{parsed}
	return nil
}
