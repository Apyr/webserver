package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func sendStatus(writer http.ResponseWriter, status int) {
	text := fmt.Sprintf("%d %s", status, http.StatusText(status))
	http.Error(writer, text, status)
}

func sendJSON(writer http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		sendStatus(writer, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write(data)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
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
