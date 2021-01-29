package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"strings"

	"gopkg.in/yaml.v2"
)

//internal
type ConfigImportStruct struct {
	CertsDir        string `yaml:"certsDir"`
	HttpPort        int    `yaml:"httpPort"`
	HttpsPort       int    `yaml:"httpsPort"`
	RedirectToHttps *bool  `yaml:"redirectToHttps"`
	Services        yaml.MapSlice
}

//internal
type ServiceImportStruct struct {
	Enabled   *bool
	Endpoints []EndpointImportStruct
}

func (s ServiceImportStruct) to(dir string) (Service, error) {
	enabled := true
	if s.Enabled != nil {
		enabled = *s.Enabled
	}
	endpoints := make([]Endpoint, 0, len(s.Endpoints))
	for _, endpoint := range s.Endpoints {
		e, err := endpoint.to(dir)
		if err != nil {
			return Service{}, err
		}
		endpoints = append(endpoints, e...)
	}
	return Service{"", enabled, endpoints}, nil
}

//internal
type EndpointImportStruct struct {
	Host         interface{}
	Path         *string
	ReverseProxy *ReverseProxyImportStruct `yaml:"reverseProxy"`
	Static       *StaticImportStruct
	Redirect     *string
}

func (e EndpointImportStruct) check() error {
	count := 0
	if e.ReverseProxy != nil {
		count++
	}
	if e.Static != nil {
		count++
	}
	if e.Redirect != nil {
		count++
	}
	if count > 1 {
		return fmt.Errorf("Too many actions in the endpoint")
	}
	if count == 0 {
		return fmt.Errorf("No actions in the endpoint %#v", e)
	}
	return nil
}

func (e EndpointImportStruct) to(dir string) ([]Endpoint, error) {
	var action Action = nil

	err := e.check()
	if err != nil {
		return nil, err
	}

	if e.ReverseProxy != nil {
		action = e.ReverseProxy.to()
	} else if e.Static != nil {
		static := e.Static.to()
		if !filepath.IsAbs(static.Dir) {
			static.Dir = filepath.Join(dir, static.Dir)
		}
		action = static
	} else {
		action = Redirect{*e.Redirect}
	}

	path := "/"
	if e.Path != nil {
		path = *e.Path
	}

	var endpoints []Endpoint
	switch host := e.Host.(type) {
	case string:
		endpoints = []Endpoint{
			Endpoint{strings.TrimSpace(host), path, action},
		}
	case []interface{}:
		endpoints = []Endpoint{}
		for _, h := range host {
			host, ok := h.(string)
			if !ok {
				return nil, fmt.Errorf("Host must be string or array of strings")
			}
			endpoints = append(
				endpoints,
				Endpoint{strings.TrimSpace(host), path, action},
			)
		}
	default:
		return nil, fmt.Errorf("Host must be string or array of strings")
	}

	return endpoints, nil
}

//internal
type ReverseProxyImportStruct struct {
	Url     string
	Replace *string
}

func (p ReverseProxyImportStruct) to() ReverseProxy {
	return ReverseProxy{p.Url, p.Replace}
}

//internal
type StaticImportStruct struct {
	Dir     string
	Page404 *string
}

func (s StaticImportStruct) to() ServeStatic {
	page := "404.html"
	if s.Page404 != nil {
		page = *s.Page404
	}
	return ServeStatic{
		Dir:     s.Dir,
		Page404: page,
	}
}

func loadFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func loadService(fileName string) (Service, error) {
	service := Service{}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return service, err
	}
	data := ServiceImportStruct{}
	err = yaml.UnmarshalStrict(content, &data)
	if err != nil {
		return service, err
	}

	dir := filepath.Dir(fileName)
	return data.to(dir)
}

func LoadConfig(fileName string) (Config, error) {
	cfg := Config{}
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return cfg, err
	}
	dir := filepath.Dir(fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cfg, err
	}
	data := &ConfigImportStruct{}
	err = yaml.UnmarshalStrict(content, &data)
	if err != nil {
		return cfg, err
	}
	redirect := true
	if data.RedirectToHttps != nil {
		redirect = *data.RedirectToHttps
	}
	if data.CertsDir == "" {
		data.CertsDir = "certs"
	}
	cfg = Config{
		CertsDir:        data.CertsDir,
		HttpPort:        data.HttpPort,
		HttpsPort:       data.HttpsPort,
		RedirectToHttps: redirect,
		Services:        []Service{},
	}
	for _, item := range data.Services {
		key, ok := item.Key.(string)
		if !ok {
			return cfg, fmt.Errorf("Service key must be string")
		}
		value, ok := item.Value.(string)
		if !ok {
			return cfg, fmt.Errorf("Service value of key \"%s\" must be string", key)
		}
		if !filepath.IsAbs(value) {
			value = filepath.Join(dir, value)
		}

		service, err := loadService(value)
		if err != nil {
			return cfg, err
		}
		service.Name = key
		cfg.Services = append(cfg.Services, service)
	}
	return cfg, nil
}
