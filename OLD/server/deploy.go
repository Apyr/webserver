package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"webserver/config"
)

type deployHandler struct {
	config.Deploy
}

func parseToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	header = strings.TrimPrefix(header, "Bearer")
	return strings.TrimSpace(header)
}

func run(env map[string]string, dir string, command []string) (string, int) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	code := 0
	result, err := cmd.CombinedOutput()
	if err != nil {
		exErr, ok := err.(*exec.ExitError)
		if ok {
			code = exErr.ExitCode()
		} else {
			return err.Error(), 1
		}
	}
	return string(result), code
}

func (deploy deployHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//auth
	authHeader := r.Header.Get("Authorization")
	token := parseToken(authHeader)
	if token == "" || token != deploy.Token {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("403 Forbidden"))
		return
	}

	//read env
	contentType := r.Header.Get("Content-Type")
	env := make(map[string]string)
	if strings.Contains(contentType, "application/json") {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(""))
			return
		}
		err = json.Unmarshal(data, &env)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("JSON parsing error: " + err.Error()))
			return
		}
	}

	//run
	output, code := run(env, deploy.Dir, deploy.Command)
	resultData := make(map[string]interface{})
	resultData["code"] = code
	resultData["output"] = output

	//send result
	data, err := json.Marshal(resultData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
