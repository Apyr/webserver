package config

import (
	"net/http"
)

type ReverseProxy struct {
	Host string
	Port int
}

func (_ ReverseProxy) actionMark() {}

type ServeStatic struct {
	Dir     string
	Page404 string
}

func (_ ServeStatic) actionMark() {}

type Redirect struct {
	To string
}

func (_ Redirect) actionMark() {}

type Action interface {
	http.Handler
	actionMark()
}

type Endpoint struct {
	Host   string
	Path   string
	Action Action
}

type Service struct {
	Enabled   bool
	Endpoints []Endpoint
}

type Config struct {
	CertsDir string
	Services []Service
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

func (cfg Config) GetHandler() http.Handler {
	hosts := make(map[string]*http.ServeMux)
	for _, service := range cfg.Services {
		if !service.Enabled {
			continue
		}
		for _, endpoint := range service.Endpoints {
			mux := hosts[endpoint.Host]
			if mux == nil {
				mux = http.NewServeMux()
				hosts[endpoint.Host] = mux
			}
			mux.Handle(endpoint.Path, endpoint.Action)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := hosts[r.Host]
		if mux != nil {
			mux.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

func DefaultConfig() Config {
	return Config{
		CertsDir: "./certs",
		Services: []Service{
			Service{
				Enabled: true,
				Endpoints: []Endpoint{
					Endpoint{
						Host: "kakotkin.ru",
						Path: "/",
						Action: ServeStatic{
							Dir:     "./static",
							Page404: "404.html",
						},
					},
				},
			},
		},
	}
}