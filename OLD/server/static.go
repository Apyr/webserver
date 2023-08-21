package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"webserver/config"
)

type staticHandler struct {
	config.Static
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func (static staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := filepath.Join(static.Dir, r.URL.Path)
	name = filepath.Clean(name)
	dir, err := isDirectory(name)
	if dir {
		name = filepath.Join(name, "index.html")
		dir, err = isDirectory(name)
	}
	if err != nil || !strings.HasPrefix(name, static.Dir) {
		name = filepath.Join(static.Dir, static.Page404)
	}
	http.ServeFile(w, r, name)
}
