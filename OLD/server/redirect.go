package server

import (
	"net/http"
	"webserver/config"
)

type redirectHandler struct {
	config.Redirect
}

func (redirect redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, redirect.To, http.StatusSeeOther)
}
