package handlers

import (
	"net/http"
	"webserver/config"
)

func NewEndpointHandler(endpoint config.Endpoint, isHTTPS bool) http.Handler {
	if !isHTTPS && *endpoint.RedirectToHTTPS {
		return RedirectToHTTPS{}
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
