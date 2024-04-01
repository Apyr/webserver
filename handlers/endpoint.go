package handlers

import (
	"net/http"
	"strings"

	"github.com/apyr/webserver/config"
)

func newEndpointHandler(endpoint config.Endpoint, isHTTPS bool) http.Handler {
	if !isHTTPS && endpoint.HTTPS != "" && endpoint.RedirectToHTTPS {
		return newRedirectToHTTPSHandler()
	}

	if endpoint.Redirect != "" {
		return newRedirectHandler(endpoint.Redirect)
	}
	if endpoint.Static != nil {
		return newStaticHandler(*endpoint.Static)
	}
	if endpoint.Proxy != nil {
		return newProxyHandler(*endpoint.Proxy)
	}
	if endpoint.RunCommand != nil {
		return newRunCommandHandler(*endpoint.RunCommand)
	}

	panic("unknown action")
}

func NewEndpointsHandler(endpoints []config.Endpoint, isHTTPS bool) http.Handler {
	mux := http.NewServeMux()
	for _, endpoint := range endpoints {
		if endpoint.Enabled != nil && !*endpoint.Enabled {
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
