package server

import (
	"net/http"
	"webserver/config"
)

func BuildHandler(cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	for _, service := range cfg.Services {
		if !service.Enabled {
			continue
		}
		for _, endpoint := range service.Endpoints {
			var handler http.Handler
			switch action := endpoint.Action.(type) {
			case config.Static:
				handler = staticHandler{action}
			case config.Redirect:
				handler = redirectHandler{action}
			case config.ReverseProxy:
				handler = reverseProxyHandler{action, endpoint.Path, nil}
			case config.Deploy:
				handler = deployHandler{action}
			default:
				panic("Porgramming error")
			}
			mux.Handle(endpoint.Host+endpoint.Path, handler)
		}
	}
	if cfg.Logging {
		return logger(mux)
	}
	return mux
}
