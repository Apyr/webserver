package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendText(writer http.ResponseWriter, status int, text string) {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Write([]byte(text))
}

func sendStatus(writer http.ResponseWriter, status int) {
	sendText(writer, status, fmt.Sprintf("%d %s", status, http.StatusText(status)))
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
