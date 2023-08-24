package handlers

import "net/http"

type redirectHandler struct {
	URL string
}

func (redirect redirectHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	http.Redirect(writer, req, redirect.URL, http.StatusSeeOther)
}
