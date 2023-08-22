package handlers

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

func (static staticHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	name := filepath.Join(static.Dir, req.URL.Path)
	name = filepath.Clean(name)
	dir, err := isDirectory(name)
	if dir {
		name = filepath.Join(name, static.Index)
	}
	if err != nil || !strings.HasPrefix(name, static.Dir) {
		name = filepath.Join(static.Dir, static.NotFound)
	}
	http.ServeFile(writer, req, name)
}
