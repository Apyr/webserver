package handlers

import (
	"net/http"
	"os/exec"
	"strings"
	"webserver/config"
)

type runCommandHandler struct {
	config.RunCommand
}

type runCommandResponse struct {
	Output string `json:"output"`
	Code   int    `json:"code"`
}

func (run runCommandHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	//auth
	authHeader := req.Header.Get("Authorization")
	token := parseToken(authHeader)
	if token == "" || token != run.Token {
		sendStatus(writer, http.StatusForbidden)
		return
	}

	//run
	output, code := execCmd(run.Command)
	sendJSON(writer, http.StatusOK, runCommandResponse{output, code})
}

func parseToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	header = strings.TrimPrefix(header, "Bearer")
	return strings.TrimSpace(header)
}

func execCmd(command []string) (string, int) {
	cmd := exec.Command(command[0], command[1:]...)

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
