package config

import (
	"net/http"
	"os"
	"path/filepath"
)

func (redirect Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, redirect.To, http.StatusSeeOther)
}

func (proxy ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := context.getProxy(proxy)
	handler.ServeHTTP(w, r)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func (static ServeStatic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := filepath.Join(static.Dir, r.URL.Path)
	dir, err := isDirectory(name)
	if dir {
		name = filepath.Join(name, "index.html")
		dir, err = isDirectory(name)
	}
	if err != nil {
		name = filepath.Join(static.Dir, static.Page404)
	}
	http.ServeFile(w, r, name)
}