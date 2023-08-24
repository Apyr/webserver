package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
