package handlers

import (
	"net/http"
	"net/url"
	"strings"
)

type Router struct {
	routes map[string][]route
}

type route struct {
	path    string
	handler http.Handler
}

func NewRouter() Router {
	return Router{make(map[string][]route)}
}

func (router Router) Add(url *url.URL, handler http.Handler) {
	host := url.Hostname()
	router.routes[host] = append(router.routes[host], route{
		path:    url.Path,
		handler: handler,
	})
}

func (router Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	host := req.URL.Hostname()

	for _, route := range router.routes[host] {
		if strings.HasPrefix(req.URL.Path, route.path) {
			route.handler.ServeHTTP(writer, req)
			return
		}
	}

	for _, route := range router.routes["*"] {
		if strings.HasPrefix(req.URL.Path, route.path) {
			route.handler.ServeHTTP(writer, req)
			return
		}
	}

	sendStatus(writer, http.StatusNotFound)
}
