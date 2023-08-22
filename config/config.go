package config

import (
	"errors"
	"net/url"
	"strings"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Ports     Ports      `yaml:"ports"`
		LogFile   string     `yaml:"logFile"`
		CertsFile string     `yaml:"certsFile"`
		LogLevel  string     `yaml:"logLevel"`
		Endpoints []Endpoint `yaml:"endpoints"`
	}

	Ports struct {
		HTTP  int `yaml:"http"`
		HTTPS int `yaml:"https"`
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
		Dir      string `yaml:"dir"`
		Index    string `yaml:"index"`
		NotFound string `yaml:"notFound"`
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
	defaultHttpPort  = 80
	defaultHttpsPort = 443
	defaultLogLevel  = "debug"
	defaultLogFile   = "webserver.log"
	defaultCertsFile = "certs.json"
	defaultNotFound  = "404.html"
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
	if config.Ports.HTTP == 0 {
		config.Ports.HTTP = defaultHttpPort
	}
	if config.Ports.HTTPS == 0 {
		config.Ports.HTTPS = defaultHttpsPort
	}
	if level := config.GetLogLevel(); level == zapcore.InvalidLevel {
		return errors.New("invalid log level")
	}

	for i := range config.Endpoints {
		if err := config.Endpoints[i].init(); err != nil {
			return err
		}
	}

	return nil
}

func (config *Config) GetLogLevel() zapcore.Level {
	switch config.LogLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.PanicLevel
	default:
		return zapcore.InvalidLevel
	}
}

func (endpoint *Endpoint) init() error {
	if endpoint.HTTPS == nil {
		val := true
		endpoint.HTTPS = &val
	}
	if endpoint.RedirectToHTTPS == nil {
		val := true
		endpoint.RedirectToHTTPS = &val
	}

	return endpoint.initActions()
}

func (endpoint *Endpoint) initActions() error {
	count := 0

	if endpoint.Redirect != "" {
		count++
	}
	if endpoint.Static != nil {
		count++
		endpoint.Static.init()
	}
	if endpoint.Proxy != nil {
		count++
	}
	if endpoint.RunCommand != nil {
		count++
		if err := endpoint.RunCommand.init(); err != nil {
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
	if static.NotFound == "" {
		static.NotFound = defaultNotFound
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

func (static *Static) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dir string
	if err := unmarshal(&dir); err == nil {
		*static = Static{Dir: dir}
		return nil
	}

	var st struct {
		Dir      string `yaml:"dir"`
		Index    string `yaml:"index"`
		NotFound string `yaml:"notFound"`
	}
	if err := unmarshal(&st); err == nil {
		*static = st
		return nil
	}

	return errors.New("can not unmarshal static")
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

	if !strings.HasPrefix(u, "http") {
		u = "http://" + u
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}

	*val = URL{parsed}
	return nil
}
