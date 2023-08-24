package handlers

import (
	"net/http"
	"net/http/httputil"
	"webserver/config"
)

type proxyHandler struct {
	config.Proxy
	handler http.Handler
}

func (proxy proxyHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if proxy.handler == nil {
		proxy.handler = httputil.NewSingleHostReverseProxy(proxy.URL.URL)
		if proxy.RemovePrefix != "" {
			proxy.handler = http.StripPrefix(proxy.RemovePrefix, proxy.handler)
		}
	}

	proxy.handler.ServeHTTP(writer, req)
}
