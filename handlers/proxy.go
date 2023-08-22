package handlers

import (
	"net/http"
	"net/http/httputil"
	"strings"
	"webserver/config"
)

type proxyHandler struct {
	config.Proxy
	handler *httputil.ReverseProxy
}

func (proxy proxyHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if proxy.handler == nil {
		proxy.handler = httputil.NewSingleHostReverseProxy(proxy.URL.URL)
	}

	if proxy.RemovePrefix != "" {
		if strings.HasPrefix(req.URL.Path, proxy.RemovePrefix) {
			req.URL.Path = strings.Replace(req.URL.Path, proxy.RemovePrefix, "", 1)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
	}

	proxy.handler.ServeHTTP(writer, req)
}
