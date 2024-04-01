package handlers

import (
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"strings"

	"github.com/apyr/webserver/config"
)

func newRedirectHandler(url string) http.Handler {
	return http.RedirectHandler(url, http.StatusSeeOther)
}

func newRedirectToHTTPSHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		url := *req.URL
		url.Scheme = "https"
		url.Host = req.Host
		http.Redirect(writer, req, url.String(), http.StatusMovedPermanently)
	})
}

func newProxyHandler(proxy config.Proxy) http.Handler {
	var handler http.Handler = httputil.NewSingleHostReverseProxy(proxy.URL.URL)
	if proxy.RemovePrefix != "" {
		handler = http.StripPrefix(proxy.RemovePrefix, handler)
	}
	return handler
}

func newStaticHandler(static config.Static) http.Handler {
	static.Dir = filepath.Clean(static.Dir)
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		name := filepath.Join(static.Dir, req.URL.Path)
		dir, err := isDirectory(name)
		if dir {
			name = filepath.Join(name, static.Index)
		}
		if err != nil || !strings.HasPrefix(name, static.Dir) {
			name = filepath.Join(static.Dir, static.NotFound)
		}
		http.ServeFile(writer, req, name)
	})
}

func newRunCommandHandler(run config.RunCommand) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		// auth
		token := req.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token != run.Token {
			sendStatus(writer, http.StatusForbidden)
			return
		}
		// run
		output, code := execCmd(run.Command)
		var response = struct {
			Output string `json:"output"`
			Code   int    `json:"code"`
		}{output, code}
		sendJSON(writer, http.StatusOK, response)
	})
}
