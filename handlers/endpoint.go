package handlers

import (
	"net/http"
	"strings"
	"webserver/config"
)

func newEndpointHandler(endpoint config.Endpoint, isHTTPS bool) http.Handler {
	if !isHTTPS && *endpoint.HTTPS && *endpoint.RedirectToHTTPS {
		return redirectToHTTPS{}
	}

	if endpoint.Redirect != "" {
		return redirectHandler{endpoint.Redirect}
	}
	if endpoint.Static != nil {
		return staticHandler{*endpoint.Static}
	}
	if endpoint.Proxy != nil {
		return proxyHandler{*endpoint.Proxy, nil}
	}
	if endpoint.RunCommand != nil {
		return runCommandHandler{*endpoint.RunCommand}
	}

	panic("unknown action")
}

func NewEndpointsHandler(endpoints []config.Endpoint, isHTTPS bool) http.Handler {
	mux := http.NewServeMux()
	for _, endpoint := range endpoints {
		if endpoint.Disabled {
			continue
		}
		handler := newEndpointHandler(endpoint, isHTTPS)
		pattern := endpoint.URL.Hostname()
		if !strings.HasPrefix(endpoint.URL.Path, "/") {
			pattern += "/"
		}
		pattern += endpoint.URL.Path
		mux.Handle(pattern, handler)
	}
	return mux
}
