package config

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type ReverseProxy struct {
	Url     string
	Replace *string
}

type Static struct {
	Dir     string
	Page404 string
}

type Redirect struct {
	To string
}

type Action interface {
	action()
}

func (_ Redirect) action()     {}
func (_ Static) action()       {}
func (_ ReverseProxy) action() {}

type Endpoint struct {
	Host   string
	Path   string
	Action Action
}

type Service struct {
	Name      string
	Enabled   bool
	Endpoints []Endpoint
}

type Config struct {
	HttpPort        int
	HttpsPort       int
	RedirectToHttps bool
	CertsDir        string
	Services        []Service
	ConfigFiles     []string
}

func (cfg Config) HttpEnabled() bool {
	return cfg.HttpPort != 0
}

func (cfg Config) HttpsEnabled() bool {
	return cfg.HttpsPort != 0
}

func (cfg Config) GetHosts() []string {
	hosts := make(map[string]bool)
	for _, service := range cfg.Services {
		if !service.Enabled {
			continue
		}
		for _, endpoint := range service.Endpoints {
			hosts[endpoint.Host] = true
		}
	}

	hostnames := []string{}
	for host := range hosts {
		hostnames = append(hostnames, host)
	}
	return hostnames
}

func (cfg Config) AsYaml() string {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (cfg Config) Watch() bool {
	watcher, err := fsnotify.NewWatcher()
	showErr := func(err error) {
		log.Printf("Watcher error: %s\n", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		showErr(err)
		return true
	}
	defer watcher.Close()

	for _, file := range cfg.ConfigFiles {
		if err := watcher.Add(file); err != nil {
			showErr(err)
			return true
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		return false
	case <-watcher.Events:
		log.Println("Config changed. Reloading...")
		return true
	case err := <-watcher.Errors:
		showErr(err)
		return true
	}
}
