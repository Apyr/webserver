package handlers

import "net/http"

type RedirectToHTTPS struct{}

func (RedirectToHTTPS) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	url := *req.URL
	url.Scheme = "https"
	url.Host = req.Host
	http.Redirect(writer, req, url.String(), http.StatusMovedPermanently)
}
