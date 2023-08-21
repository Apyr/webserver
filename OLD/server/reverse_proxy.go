package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"webserver/config"
)

type reverseProxyHandler struct {
	config.ReverseProxy
	path    string
	handler *httputil.ReverseProxy
}

func (proxy reverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if proxy.handler == nil {
		u, err := url.Parse(proxy.Url)
		if err != nil {
			log.Fatalf("Reverse proxy url parsing error: %s", err)
		}
		proxy.handler = httputil.NewSingleHostReverseProxy(u)
	}
	if proxy.Replace != nil {
		if strings.HasPrefix(r.URL.Path, proxy.path) {
			r.URL.Path = strings.Replace(r.URL.Path, proxy.path, *proxy.Replace, 1)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
		}
	}
	proxy.handler.ServeHTTP(w, r)
}
